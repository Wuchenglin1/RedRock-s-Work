package main

import (
	"fmt"
	"sync"
	"time"
)

func juan (){
	defer wg.Done()
	newTicker := time.NewTicker(24 * time.Hour)
	if time.Now().Hour()==2{
		for t := range newTicker.C{
			fmt.Println(t)
			fmt.Println("谁能比我卷")
		}
	}
}

func zao8(){
	defer wg.Done()
	newTicker := time.NewTicker(24 * time.Hour)
	if time.Now().Hour()==6{
		for t := range newTicker.C{
			fmt.Println(t)
			fmt.Println("早八算什么，早六才是吾辈应起之时")
		}
	}
}

func wuhu(){
	defer wg.Done()
	newTicker := time.NewTicker(time.Second * 30)
	for t := range newTicker.C{
		fmt.Println(t)
		fmt.Println("芜湖！起飞！")
	}
}

var wg sync.WaitGroup

func main(){
	wg.Add(3)
	go juan()
	go zao8()
	go wuhu()
	wg.Wait()
}