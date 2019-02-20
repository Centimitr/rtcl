package rtcl

func NewRTCLFromFile(filename string) (rtcl *RTCL, err error) {
	ast, err := ParseFile(filename)
	if err != nil {
		return
	}
	return NewRTCLFromAST(ast)
}

func NewRTCLFromAST(ast *node) (r *RTCL, err error) {
	cur := newASTCursor(ast)

	meta := &Meta{Attributes: make(map[string]string)}
	r = &RTCL{Meta: meta,}

	if cur.locateFromRoot("article.meta.Args", 3) {
		for _, child := range cur.ptr.children() {
			meta.addArg(child.val)
		}
	}

	if cur.locateFromRoot("article.meta.kvs", 3) {
		for _, node := range cur.ptr.children() {
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

	if cur.locateFromRoot("article.content", 2) {
		wrapper := cur.ptr.child
		HandleBlock(wrapper)
		r.Content = wrapper.representation
	}

	return
}

func HandleBlock(node *node) {
	if node == nil {
		return
	}

	handleChildren := func() {
		for _, child := range astChildren(node) {
			HandleBlock(child)
		}
	}

	h := DefaultHandlers.Match(node)

	if h.Handle != nil {
		h.Handle(node, handleChildren)
	}

	//fmt.Println(node.typ, node.representation)
}
