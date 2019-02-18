package main

import "sync"

func NewRTCLFromFile(filename string) (rtcl *Container, err error) {
	ast, err := ParseFile(filename)
	if err != nil {
		return
	}
	return NewRTCLFromAST(ast)
}

func NewRTCLFromAST(ast *node) (r *Container, err error) {

	meta := &Meta{Attributes: make(map[string]string)}
	content := &Content{&Block{}}
	r = &Container{Meta: meta, Content: content,}

	if ast.locateFromRoot("article.meta.args", 3) {
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
		HandleBlock(ast.ptr.child)
	}

	return
}

func HandleBlock(node *node) {
	if node == nil {
		return
	}
	h := Handlers.Match(node)
	var processChildren = true

	if h.Pre != nil {
		processChildren = h.Pre(node)
	}

	if processChildren {
		wg := &sync.WaitGroup{}
		for _, child := range astChildren(node) {
			wg.Add(1)
			go func() {
				HandleBlock(child)
				wg.Done()
			}()
		}
		wg.Wait()
	}

	if h.Post != nil {
		h.Post(node)
	}
}
