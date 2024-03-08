package main

import (
	"fmt"
	"runtime"
	"time"
)

// 利用runtime包来实现
func a() {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, s := range str {
		fmt.Printf("%c", s)
		runtime.Gosched()
	}
}

func b() {
	str := "012345678910111213141516171819202122232425"
	for i := 0; i < len(str); i++ {
		if i >= 10 {
			fmt.Printf("%s", str[i:i+2])
			i++
		} else {
			fmt.Printf("%c", str[i])
		}
		runtime.Gosched()
	}
}

func main() {
	go a()
	go b()
	time.Sleep(time.Second)
}
