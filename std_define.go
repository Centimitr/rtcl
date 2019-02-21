package rtcl

type Section struct {
	*Container
	Name string
}

type Define struct {
	Dict map[string]string
}

type Text struct {
	String string
}

type Paragraph struct {
	Fragments []interface{}
	String    string
}

func (p *Paragraph) UpdateString() {
	var s string
	sep := " "
	for _, frag := range p.Fragments {
		switch v := frag.(type) {
		case *Text:
			if s != "" {
				s += sep
			}
			s += v.String
		}
	}
	p.String = s
}

type List struct {
	Options string
}

type TaskList struct {
	Tasks interface{}
}
