package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

func Fn[T any]() {
	fmt.Printf("%T\n", *new(T))
}

func TestType(t *testing.T) {
	Fn[int]()
	Fn[rune]()
	Fn[byte]()
}
