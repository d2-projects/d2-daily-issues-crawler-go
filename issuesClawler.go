package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/robfig/cron"
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

func run() {
	fmt.Println("Ready! Gooo! %v", time.Now())
	nameString := dayString() //dd
	dateString := datString() //yyyy.mm.dd
	filename := nameString + ".md"
	//go文件要在根目录，todo 判断今天的文件是否已经存在
	//var dir string = "201810/01.md"
	mkdir4month(dateString)
	//create markdown file like 01.md
	createMarkDown(dateString, filename)
	//start scrap
	mdContext := scrape()
	fmt.Println(mdContext)
	//keep write
	writeMdContext(mdContext, filename)
	//退出目录
	exitDir()
	//todo git
	gitPushDaily(nameString)
}


func scrape() string {
	//var mdContext string  //
	mapList := make(map[string]string)
	response := getResponse(basUrl+issuesUrl)
	// 获取issue主页
	dom, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatalf("失败原因", response.StatusCode)
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
	mdContext := map2string(mapList)
	return mdContext
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
			name := checkNameNode(lineContext[2])
			node := checkNameNode(lineContext[3])
			m[key1] = "\n" + name +"\n"+node+"\n"+lineContext[4]
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
func mkdir4month(date string) string{
	dateArray := strings.Split(date,".")
	dirnam := dateArray[0] + dateArray[1]
	day := dateArray[2]
	if day == "01" {
		cmd := exec.Command("mkdir",dirnam)
		out, err := cmd.Output()
		if err != nil {
			println(err.Error())
			return dirnam
		}
		print(string(out))
	}
	return dirnam
}

func exitDir() {
	cmd := exec.Command("cd","..")
	out, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return
	}
	print(string(out))
	cmd2 := exec.Command("pwd")
	out2, err2 := cmd2.Output()
	if err2 != nil {
		println(err2.Error())
		return
		print(string(out2))
	}
}

func getResponse(url string) *http.Response {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0")
	response, _ := client.Do(request)
	return response
}

func dayString() string {
	y, m, d := time.Now().Date()
	mStr := fmt.Sprintf("%d", m)
	dStr := fmt.Sprintf("%d", d)
	yStr := fmt.Sprintf("%d", y)
	if m < 10 {
		mStr = fmt.Sprintf("0%d", m)
	}
	if d < 10 {
		dStr = fmt.Sprintf("0%d", d)
	}
	fmt.Sprintf("%d-%s", yStr,mStr)
	return fmt.Sprintf("%d%s%s", y, mStr,dStr)

}

func datString() string {
	y, m, d := time.Now().Date()
	mStr := fmt.Sprintf("%d", m)
	dStr := fmt.Sprintf("%d", d)
	if m < 10 {
		mStr = fmt.Sprintf("0%d", m)
	}
	if d < 10 {
		dStr = fmt.Sprintf("0%d", d)
	}
	return fmt.Sprintf("%d.%s.%s", y, mStr, dStr)

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

func createMarkDown(date string, filename string) {
	// open output file
	fo, err := os.Create(filename)
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
editor:
- name: FairyEver
  url: https://github.com/FairyEver
- name: ishenyi
  url: https://github.com/ishenyi
- name: Jiiiiiin
  url: https://github.com/Jiiiiiin
- name: sunhaoxiang
  url: https://github.com/sunhaoxiang
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
