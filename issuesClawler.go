package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	basUrl string = "https://github.com"
	issuesUrl string = "/d2-projects/d2-awesome/issues"
)

func main() {

	nameString := dayString()
	dateString := datString()
	filename := nameString + ".md"

	//go文件要在根目录，判断今天的文件是否已经存在
	//todo 自动创建文件夹，创建目录，修改readme
	//var dir string = "site/daily/post/2018/10/17.md"

	//create markdown file
	createMarkDown(dateString, filename)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

  	var slide1 string = `- name: 新闻
  list:`
	var slide2 string = `
- name: 开源项目
  list:`
	var slide3 string = `
- name: 分享
  list:`
	var slide4 string = `
- name: 教程
  list:`
	var slide5 string = `
- name: 工具
  list:`
	var slide6 string = `
- name: 招聘
  list:`
	var slide7 string = `
- name: 设计
  list:`

	var num1 int = 0
	var num2 int = 0
	var num3 int = 0
	var num4 int = 0
	var num5 int = 0
	var num6 int = 0
	var num7 int = 0

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
			// 进入issue子页面获取内容
			href = strings.TrimSpace(href)
			res := getResponse(basUrl+href)
			childDom, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatalf("子页面失败原因", response.StatusCode)
			}
			childDom.Find("pre").Each(func(i int, s *goquery.Selection){

				//直接拼接不可
				//mdContext += s.Text()
				//fmt.Println(mdContext)

				// 写到这里只是简单的对 issues 的内容作了一个拼接，下面要解决的问题，就是 MapReduce 的问题。
				//参考 https://github.com/happyer/distributed-computing/blob/master/src/mapreduce/README.md
				//我准备用最简的办法来做😆，很不优雅
				//strings.trim函数这里对多行的string存在bug，必须做2次截取
				//issueString := s.Text()

				//2018-11-21 10:32:52 strings.trim处理多行字符串存在bug，重新调整一种方法
				lineContext := strings.Split(s.Text(), "\n")
				//fmt.Println(lineContext)
				if len(lineContext) > 5  {
					typeArray := strings.Split(lineContext[0], ": ")
					if len(typeArray) >= 2 {
						if typeArray[1] == "新闻" {
							num1 += 1
							slide1 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
						if typeArray[1] == "开源项目" {
							num2 += 1
							slide2 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
						if typeArray[1] == "分享" {
							num3 += 1
							slide3 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
						if typeArray[1] == "教程" {
							num4 += 1
							slide4 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
						if typeArray[1] == "工具" {
							num5 += 1
							slide5 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
						if typeArray[1] == "招聘" {
							num6 += 1
							slide6 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
						if typeArray[1] == "设计" {
							num7 += 1
							slide7 += "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
						}
					}
				}

			})
		}

	})
	var mdContext string =""
	if num1>0 {
		mdContext += slide1
	}
	if num2>0 {
		mdContext += slide2
	}
	if num3>0 {
		mdContext += slide3
	}
	if num4>0 {
		mdContext += slide4
	}
	if num5>0 {
		mdContext += slide5
	}
	if num6>0 {
		mdContext += slide6
	}
	if num7>0 {
		mdContext += slide7
	}
	fmt.Println(mdContext)

	var foot string = `
---

<daily-list v-bind="$page.frontmatter"/>`
	// 写进文件
	if _, err = f.WriteString(mdContext + foot); err != nil {
		println(err.Error())
		panic(err)
	}


}

/**
* 返回response
*/
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
	return fmt.Sprintf("%s", dStr)

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
