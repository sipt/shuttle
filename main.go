package main

import "fmt"

type A struct {
	Name string
}

func main() {
	s := []A{
		{Name: "123"},
		{Name: "asd"},
	}
	d := make([]A, len(s))
	copy(d, s)
	set(d)
	fmt.Println(s[0])
}

func set(s []A) {
	s[0].Name = "aaa"
}
