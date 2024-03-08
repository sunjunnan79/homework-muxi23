package main

import (
	"fmt"
)

type Builder[T any] struct {
	data []T
}

// NewBuilder 初始化的
func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		data: make([]T, 0),
	}
}

// Add 添加项的
func (b *Builder[T]) Add(value T) {
	b.data = append(b.data, value)
}

// PrintBuilder 输出的
func (b *Builder[T]) PrintBuilder() {
	fmt.Println(b.data)
}

// GetString 获取字符串类型的
func (b *Builder[T]) GetString() (s string) {
	s = "["
	for i, j := range b.data {
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%v", j)
	}
	s += "]"
	return
}

// Len 计算长度的
func (b *Builder[T]) Len() int {
	return len(b.data)
}

// ReSet 清除缓存的
func (b *Builder[T]) ReSet() {
	b.data = make([]T, 0)
}

func main() {
	intBuilder := NewBuilder[int]()
	intBuilder.Add(10)
	intBuilder.Add(20)
	intBuilder.Add(30)
	fmt.Printf("长度:%d\n", intBuilder.Len())
	intBuilder.PrintBuilder()
	fmt.Println("内容:" + intBuilder.GetString())

	intBuilder.ReSet()
	intBuilder.Add(30)
	intBuilder.Add(20)
	intBuilder.Add(10)
	intBuilder.PrintBuilder()

	strBuilder := NewBuilder[string]()
	strBuilder.Add("hello")
	strBuilder.Add("World!")
	strBuilder.PrintBuilder()
}
