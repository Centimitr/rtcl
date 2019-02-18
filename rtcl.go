package main

type Container struct {
	Meta    *Meta
	Content *Content
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

type Content struct {
	Block *Block
}

type Block struct {
	Command string
	Content
}

func (rtcl *Container) Print() {
	_ = printJson(rtcl)
}
