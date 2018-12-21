package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	basUrl string = "https://github.com"
	issuesUrl string = "/d2-projects/d2-awesome/issues"
)

func main() {
	i := 0
	c := cron.New()
	spec := "0 0 12-14 * * ?"
	c.AddFunc(spec, func() {
		i++
		run()
	})
	c.Start()
	select{}
}

//test locally ,pls chang run(){} as main(){}, add dismiss git operation
//本地测试的时候，把run方法改为main方法，并注释掉commit操作，祝使用顺利！
func run() {
	fmt.Println("Ready! Gooo! %v", time.Now())
	y, m, d := time.Now().Date()
	mStr := fmt.Sprintf("%d", m)
	dStr := fmt.Sprintf("%d", d)
	dyestStr := fmt.Sprintf("%d", d-1)
	if m < 10 {
		mStr = fmt.Sprintf("0%d", m)
	}
	if d < 10 {
		dStr = fmt.Sprintf("0%d", d)
	}
	if d < 11 {
		dyestStr = fmt.Sprintf("0%d", d-1)
	}
	dateString := fmt.Sprintf("%d.%s.%s", y, mStr, dStr) //yyyy.mm.dd
	urlParam := fmt.Sprintf("%d-%s-%s", y, mStr, dyestStr) //yyyy-mm-dd
	dir := "./" + fmt.Sprintf("%d/%s", y, mStr)
	path := dir + "/" + fmt.Sprintf("%s", dStr) + ".md"
	//go文件要在根目录，判断今天的文件是否已经存在 例如 path = "./2018/11/11.md"
	mkdir4month(dir)
	//create markdown file like 01.md
	createMarkDown(dateString, path)
	//start scrap
	mdContext := scrape(urlParam)
	fmt.Println(mdContext)
	//keep write
	writeMdContext(mdContext, path)
	//git
	gitPushDaily(urlParam)
}


func scrape(urlParam string) string {
	//var mdString string = ``
	mdMpList := make(map[string]string)

	//todo if crawler cron time change pls check here for issues search param! It's 14:00:00+08:00 everyday
	urlQueryParam := "?q=created%3A>" + urlParam + "T14%3A00%3A00%2B08%3A00+is%3Aopen"
	url := issuesUrl+urlQueryParam

	for {
		mdMpList = getContextListPerPage(mdMpList,basUrl+url)
		//判断是否要获取issue下一页数据
		response := getResponse(basUrl+url)

		// 获取issue主页
		dom, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatalf("获取下一页失败原因", response.StatusCode)
		}
		nextUrl := "" //重置
		dom.Find("a[class=next_page]").Each(func(i int, selection *goquery.Selection){
			fmt.Println(selection.Attr("href"))
			nextUrl,_ = selection.Attr("href")
		})

		if nextUrl == "" {
			break
		} else {
			url = nextUrl
		}
	}
	mdString := map2string(mdMpList)
	return mdString
}

func getContextListPerPage(mapList map[string]string, url string) map[string]string {
	//mapList := make(map[string]string)
	response := getResponse(url)
	dom, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatalf("获取issues失败原因", response.StatusCode)
	}
	dom.Find("a[data-hovercard-type=issue]").Each(func(i int, selection *goquery.Selection) {
		// 获取issue 的 href
		href, IsExist := selection.Attr("href")
		fmt.Println(href)
		if IsExist == true {
			href = strings.TrimSpace(href)
			res := getResponse(basUrl+href)
			childDom, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatalf("子页面失败原因", response.StatusCode)
			}
			childDom.Find("pre").Each(func(i int, s *goquery.Selection){
				// shape s.Text into map
				mapSingleList := shapeText2Map(s.Text())
				if len(mapSingleList) != 0 {
					mapList = mapListAdds(mapSingleList,mapList)
				}
			})
		}
	})
	return mapList
}

func shapeText2Map(s string) map[string]string {
	var m map[string]string
	m = make(map[string]string)
	lineContext := strings.Split(s, "\n")
	//fmt.Println(lineContext)
	if len(lineContext) >= 5  {
		typeArray := strings.Split(lineContext[0], ": ")
		if len(typeArray) >= 2 {
			key1 := typeArray[1]
			for i:=2; i<= len(lineContext)-1; i++ {
				ConArray := strings.Split(lineContext[i], ": ")
				if len(ConArray) >= 2 {
					if strings.Contains(ConArray[0],"name") || strings.Contains(ConArray[0],"note") {
						//fix full :
						m[key1] = m[key1] + "\n" + checkNameNode(lineContext[i])
					} else {
						m[key1] = m[key1] + "\n" + lineContext[i]
					}
				}
			}
		}
	}
	return m
}

func checkNameNode(mp string) string {
	ns := strings.Replace(mp, ": ", "?:?", 1)
	fs := strings.Replace(ns, ": ", " ", -1)
	return strings.Replace(fs, "?:?", ": ", 1)
}

func mapListAdds(mapSingleList map[string]string, mapList map[string]string) map[string]string {
	if len(mapList) == 0 {
		for k1, v1 := range mapSingleList {
			mapList[k1] = v1
		}
	} else {
		for k, v := range mapSingleList {
			// 查找 key 是否存在
			if _, ok := mapList[k] ; ok {
				mapList[k] = mapList[k] + v
			} else {
				mapList[k] = v
			}
		}
	}
	return mapList
}

func map2string(mapp map[string]string) string{
	var s string = ""
	var n string = "- name: "
	var l string = `  list:`
	if len(mapp) != 0 {
		for k, v := range mapp {
			s += n + k + "\n" + l + v +"\n"
		}
		return s
	} else {
		return s
	}
}

//每月1号自动创建文件夹，创建目录；进入目录
func mkdir4month(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Println("Create Directory OK!" + path)
	}
}

func getResponse(url string) *http.Response {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0")
	response, _ := client.Do(request)
	return response
}

func writeMdContext(md string, fn string) {
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var foot string = `
---

<daily-list v-bind="$page.frontmatter"/>`
	// 写进文件
	if _, err = f.WriteString(md + foot); err != nil {
		println(err.Error())
		panic(err)
	}
}

func createMarkDown(date string, path string) {
	// open output file
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(fo)
	var title string = `---
pageClass: page-daily-list
date: {+dateString+}
title: 日报 {+dateString+}
meta:
- itemprop: name
  content: 日报 {+dateString+}
- name: description
  itemprop: description
  content: 今天的新发现
list:
`
	w.WriteString(strings.Replace(title, "{+dateString+}",date, -1) )
	w.Flush()
}

func gitPushDaily(fn string){
	gitPull()
	gitAddAll()
	gitCommit(fn)
	gitPush()
}

func gitPull() {
	app := "git"
	arg0 := "pull"
	arg1 := "origin"
	arg2 := "master"
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitAddAll() {
	app := "git"
	arg0 := "add"
	arg1 := "."
	cmd := exec.Command(app, arg0, arg1)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitCommit(date string) {
	app := "git"
	arg0 := "commit"
	arg1 := "-am"
	arg2 := date
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}
func gitPush() {
	app := "git"
	arg0 := "push"
	arg1 := "origin"
	arg2 := "master"
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}
