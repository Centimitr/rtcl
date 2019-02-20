package rtcl

import (
	"fmt"
	"time"
)

func ExampleFile() {

	now := time.Now()

	rtcl, err := NewRTCLFromFile("test.rtcl")
	check(err)

	rtcl.Print()
	//_ = rtcl
	fmt.Println(time.Since(now))
	// Output:
}
