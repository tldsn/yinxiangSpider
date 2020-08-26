package main

import (
	"fmt"
	"runtime"
	"yinxiangSpider/spider"
)

func main() {
	runtime.GOMAXPROCS(4)
	taskList, _ := spider.GetNoteUrl()
	var tasks = make(chan string)
	fmt.Println("...")
	a1 := make(chan bool)
	a2 := make(chan bool)
	a3 := make(chan bool)

	// for i := 0; i < 3; i++ {
	// 	go woker(tasks)
	// }
	go woker(tasks, a1)
	go woker(tasks, a2)
	go woker(tasks, a3)

	for _, v := range taskList {
		tasks <- v
	}
	close(tasks)
	<-a1
	fmt.Println("工作完成1")
	<-a2
	fmt.Println("工作完成2")
	<-a3
	fmt.Println("工作完成3")
}

func woker(taskList chan string, a chan bool) {
	for taskStr := range taskList {
		err := spider.EnterNoteUrl(taskStr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	a <- true
}
