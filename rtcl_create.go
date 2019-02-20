package rtcl

func NewRTCLFromFile(filename string) (rtcl *RTCL, err error) {
	ast, err := ParseFile(filename)
	if err != nil {
		return
	}
	return NewRTCLFromAST(ast)
}

func NewRTCLFromAST(ast *node) (r *RTCL, err error) {

	meta := &Meta{Attributes: make(map[string]string)}
	r = &RTCL{Meta: meta,}

	if ast.locateFromRoot("article.meta.Args", 3) {
		for _, node := range astChildren(ast.ptr) {
			meta.addArg(node.val)
		}
	}

	if ast.locateFromRoot("article.meta.kvs", 3) {
		for _, node := range astChildren(ast.ptr) {
			var k, v string
			if node.child != nil {
				k = node.child.val
			}
			if node.child.sibling != nil {
				v = node.child.sibling.val
			}
			meta.addAttribute(k, v)
		}
	}

	if ast.locateFromRoot("article.content", 2) {
		wrapper := ast.ptr.child
		HandleBlock(wrapper)
		r.Content = wrapper.representation
	}

	return
}

func HandleBlock(node *node) {
	if node == nil {
		return
	}
	h := Handlers.Match(node)

	var needProcessChildren = true

	if h.Pre != nil {
		needProcessChildren = h.Pre(node)
		h.Pre(node)
	}

	if needProcessChildren {
		for _, child := range astChildren(node) {
			HandleBlock(child)
		}
	}

	if h.Post != nil {
		h.Post(node)
	}

	//fmt.Println(node.typ, node.representation)
}
