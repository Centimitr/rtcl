package rtcl

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func printJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "    ")
	fmt.Println(string(b))
	return err
}

func printJSONInChrome(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	dataURL := strings.Replace("data:text/html,<script>console.dir(JSON.parse('$json'))</script>", "$json", string(b), 1)

	var prefixes []string
	switch runtime.GOOS {
	case "darwin":
		prefixes = []string{"open"}
	case "windows":
		prefixes = []string{"cmd", "/c", "start"}
	default:
		prefixes = []string{"xdg-open"}
	}
	args := []string{"chrome", dataURL, "--auto-open-devtools-for-tabs"}
	return exec.Command(prefixes[0], append(prefixes[1:], args...)...).Run()
}

type args struct {
	vs     []string
	first  string
	second string
	third  string
	last   string
}

func newArgs(s string) *args {
	a := &args{vs: strings.Split(s, " ")}
	a.first = a.value(0)
	a.second = a.value(0)
	a.third = a.value(0)
	a.last = a.value(len(a.vs) - 1)
	return a
}

func (a *args) value(index int) string {
	if index >= len(a.vs) {
		return ""
	}
	return a.vs[index]
}
