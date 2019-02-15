package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

func ParseFile(filename string) (ast *node, err error) {
	t := time.Now()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	fmt.Println("[Read]", time.Since(t).Nanoseconds(), "ns")
	fmt.Println()

	t = time.Now()
	ast = Parse(string(content))
	fmt.Println("[Parse]", time.Since(t).Nanoseconds(), "ns")
	return
}

func main() {

	rtcl, err := NewRTCLFromFile("test.rtcl")
	check(err)

	rtcl.Print()
}
