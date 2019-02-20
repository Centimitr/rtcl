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

type Args struct {
	Slice  []string
	First  string
	Second string
	Third  string
	Last   string
}

func NewArgsFromString(s string) *Args {
	return (&Args{Slice: strings.Split(s, " ")}).Parse()
}

func NewArgsFromSlice(vs []string) *Args {
	return (&Args{Slice: vs}).Parse()
}

func (a *Args) Parse() *Args {
	a.First = a.Value(0)
	a.Second = a.Value(1)
	a.Third = a.Value(2)
	a.Last = a.Value(len(a.Slice) - 1)
	return a
}

func (a *Args) Append(s string) *Args {
	a.Slice = append(a.Slice, s)
	a.Parse()
	return a
}

func (a *Args) Value(index int) string {
	if index >= len(a.Slice) {
		return ""
	}
	return a.Slice[index]
}
