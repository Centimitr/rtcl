package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func printJson(v interface{}) error {
	b, err := json.MarshalIndent(v, "", " ")
	fmt.Println(string(b))
	return err
}
