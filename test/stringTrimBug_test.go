package test

import (
	"fmt"
	"strings"
)

func main() {
	var a string = `- name: 开源项目
  list:
  - name: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    url: https://github.com/mafengwo/vue-drag-tree-table`
	var twoLineString string= `- name: 开源项目
  list:`
	var firstline string = "- name: 开源项目"
	var secondline string = "list:"

	fmt.Println(strings.Trim(a,twoLineString))
/*
fengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    url: https://github.com/mafengwo/vue-drag-tree-tab
*/

	fmt.Println(strings.Trim(a,firstline))
/*
list:
  - name: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    url: https://github.com/mafengwo/vue-drag-tree-tabl
*/
	b := strings.Trim(a,firstline)
	fmt.Println(strings.Trim(b,secondline))
/*
list:
	- name: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
url: https://github.com/mafengwo/vue-drag-tree-tab
*/
	ss := strings.Split(a,"\n")
	fmt.Println("\n" + ss[2]+"\n"+ss[3]+"\n"+ss[4])
	fmt.Println("-------------------")
	fmt.Println(checkNameNode(`  - name: sorrycc/awesome-f2e-libs: 🎉 整理我平时关注的前端库。`))
}

func checkNameNode(mp string) string {
	ns := strings.Replace(mp, ": ", "?:?", 1)
	fs := strings.Replace(ns, ": ", " ", -1)
	return strings.Replace(fs, "?:?", ": ", 1)
}