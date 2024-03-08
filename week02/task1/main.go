package main

import (
	"fmt"
	"time"
)

// 利用通道来实现
func a(ch1, ch2 chan int) {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, s := range str {
		<-ch1
		fmt.Printf("%c", s)
		ch2 <- 1
	}
}

func b(ch1, ch2 chan int) {
	str := "012345678910111213141516171819202122232425"
	for i := 0; i < len(str); i++ {
		<-ch2
		if i >= 10 {
			fmt.Printf("%s", str[i:i+2])
			i++
		} else {
			fmt.Printf("%c", str[i])
		}
		ch1 <- 1 // 确保在打印完内容后再接收 ch1 数据
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go a(ch1, ch2)
	go b(ch1, ch2)
	ch1 <- 1
	time.Sleep(10 * time.Second)
}
