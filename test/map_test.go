package test

import (
	"fmt"
	"strings"
)

func main() {
	capitals := map[string] string {"France":"Paris", "Italy":"Rome", "Japan":"Tokyo" }
	for key := range capitals {
		fmt.Println("Map item: Capital of", key, "is", capitals[key])
	}
	maplist := map[string]string {}
	if len(maplist) == 0 {
		fmt.Println("null")
	}else{
		fmt.Println(len(maplist))
	}

	var s string = `- name: 开源项目
  list:
  - name: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格
    url: https://github.com/mafengwo/vue-drag-tree-table`
	var m map[string]string
	m = make(map[string]string)
	lineContext := strings.Split(s, "\n")
	//fmt.Println(lineContext)
	if len(lineContext) >= 5  {
		typeArray := strings.Split(lineContext[0], ": ")
		if len(typeArray) >= 2 {
			key1 := typeArray[1]
			m[key1] = "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
		}
	}
	for k,v := range m{

		fmt.Println(k,v)
	}

	mall := make(map[string]string)
	mall["开源项目"] = ` 
  - name: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格1
    note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格1`
	mall["飞翔"] = ` 
  - name: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格2
    note: mafengwo/vue-drag-tree-table: vue 可以拖拽的树形表格2`

	ms := mapListAdds(m, mall)
	fmt.Println("------------------")
	for k,v := range ms{

		fmt.Println(k,v)
	}

	sr := map2string(ms)
	fmt.Println("------------------")
	fmt.Println(sr)
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
			m[key1] = "\n" + lineContext[2]+"\n"+lineContext[3]+"\n"+lineContext[4]
		}
	}
	return m
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