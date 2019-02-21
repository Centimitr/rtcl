package rtcl

import (
	"fmt"
	"strings"
)

func HTMLs(vs []interface{}) (html string) {
	for _, v := range vs {
		html += HTML(v)
	}
	return
}

func wrap(selector string, children ...interface{}) string {
	vs := strings.Split(selector, ".")
	tag := vs[0]

	following := ""
	if len(vs) > 1 {
		following = fmt.Sprintf(` class="%s"`, strings.Join(vs[1:], " "))
	}
	return fmt.Sprintf(`<%s%s>%s</%s>`, tag, following, HTMLs(children), tag)
}

func HTML(v interface{}) (html string) {
	switch v := v.(type) {
	case *Section:
		return wrap("section", v.Children)
	case *Paragraph:
		return wrap("p", v.String)
	case *Container:
		return HTMLs(v.Children)
	case []interface{}:
		return HTMLs(v)
	case *TaskList:
		return wrap("ul", v.Tasks)
	case *List:
		return wrap("ul", v.Options)
	case *Define:
		return wrap("section.define", "<div>DEFINE</div>", v.Dict)
	case map[string]string:
		for k, v := range v {
			html += wrap("div.item", wrap("div.key", k), wrap("div.value", v))
		}
	case string:
		html = v
	default:
		fmt.Println("!!", v)
	}
	return
}
