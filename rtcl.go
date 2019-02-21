package rtcl

import "html/template"

type RTCL struct {
	Meta    *Meta
	Content interface{}
}

func (rtcl *RTCL) Print() {
	check(printJSON(rtcl))
}

func (rtcl *RTCL) HTML() template.HTML {
	return template.HTML(HTML(rtcl.Content))
}

type Meta struct {
	Title      string
	Subtitle   string
	Arguments  []string
	Attributes map[string]string
}

func (m *Meta) addArg(s string) {
	switch len(m.Arguments) {
	case 0:
		m.Title = s
	case 1:
		m.Subtitle = s
	}
	m.Arguments = append(m.Arguments, s)
}

func (m *Meta) addAttribute(k string, v string) {
	m.Attributes[k] = v
}

func NewContainerFromNode(node *node) *Container {
	c := &Container{}
	for _, child := range astChildren(node) {
		if child.representation != nil {
			c.Append(child.representation)
		}
	}
	return c
}

type Container struct {
	Children []interface{}
}

func (c *Container) Append(v interface{}) {
	c.Children = append(c.Children, v)
}
