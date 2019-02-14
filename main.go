package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

func ParseFile(filename string) error {
	t := time.Now()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Println("[Read]", time.Since(t).Nanoseconds(), "ns")
	fmt.Println()

	t = time.Now()
	Parse(string(content))
	fmt.Println("[Parse]", time.Since(t).Nanoseconds(), "ns")
	return nil
}

func Parse(s string) {
	l := newLexer(s)
	go l.run()
	for item := range l.items {
		fmt.Println(item.typ, item.val)
	}
}

func main() {
	_ = ParseFile("test.rtcl")
}
