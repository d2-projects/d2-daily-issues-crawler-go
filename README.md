# issues-crawler-go

go 写的爬虫。 自动整理 Github issue，并归类。

主要是 D2 Awesome Daily 的自动化工具。

## How to use

1. 执行 `./issuesClawler`
2. 程序每天自动执行一次，产生文件 2018mmdd.md 在根目录
3. copy 内容到 d2-awesome , 记得增加 d2-awesome readme.md的对应日期

## To Do

1. go 文件放在根目录，可创建文件
2. ~~自动提交~~
3. ~~处理 issue 的评论~~
4. ~~增加冒号替换，解决格式问题~~

## 初始化

**已无需初始化,支持任意类型，不限于以下几类**

~~在使用 ios 的捷径分享的第一步时，选择分享的类型，目前爬虫支持的类型有：

- 新闻
- 开源项目
- 分享
- 教程
- 工具
- 招聘
- 设计
- 资源~~

## milestone

2018-11-21 23:47:03

1. 代码美化
2. 增加 git 提交
3. fix "：" 导致 vuepress 转换报错
4. 增加不完全 test 用例

2018-12-11

1. issue 爬取增加时间限制，仅限 t-1 14:00以后的 issue，避免 issue 未及时关闭导致内容重复