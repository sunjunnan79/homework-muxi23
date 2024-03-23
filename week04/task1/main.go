package main

import (
	"fmt"
	"unsafe"
)

func bytesToString(b []byte) string {
	//强制转化后直接返回一个string类型的指针
	return *(*string)(unsafe.Pointer(&b))
}

func bytesToSliceInt(b []byte) []int {
	//强制转化后直接返回一个string类型的指针
	return *(*[]int)(unsafe.Pointer(&b))
}

func main() {
	bytes := []byte("Good morning!")
	str := bytesToString(bytes)
	fmt.Println(str)
	bytes[5] = 'n'
	// 修改原始字节数组
	fmt.Println(str)

	////实现了将bytes解释为int[]再输出,想不明白输出结果,好像是因为byte占一个字节,int占8个,但输出为什么是这样不能理解
	//n := bytesToSliceInt(bytes)
	//fmt.Println(n)
}
