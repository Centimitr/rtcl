package main

type Container struct {
	Meta    *Meta
	Content *Content
}

type Meta struct {
	Args []string
	KVs  map[string]string
}

type Content struct {
	Block *Block
}

type Block struct {
	Cmd string
	Content
}

func (rtcl *Container) Print() {
	_ = printJson(rtcl)
}

func NewRTCLFromAST(ast *node) *Container {
	r := &Container{
		Meta:    &Meta{},
		Content: &Content{&Block{}},
	}

	return r
}

func NewRTCLFromFile(filename string) (rtcl *Container, err error) {
	ast, err := ParseFile(filename)
	if err != nil {
		return
	}
	rtcl = NewRTCLFromAST(ast)
	return
}
