package main

func main() {

	rtcl, err := NewRTCLFromFile("test.rtcl")
	check(err)

	rtcl.Print()
}
