package rtcl

type Handler struct {
	Match func(*node) bool
	Pre   func(*node) bool
	Post  func(*node)
}

func NewWrapperPost(name string) func(*node) {
	return func(node *node) {
		b := &Block{Command: name}
		for _, child := range astChildren(node) {
			if child.representation != nil {
				b.AddChild(child.representation)
			}
		}
		node.representation = b
	}
}

var defaultHandler = &Handler{Post: NewWrapperPost(".")}

var omitHandler = &Handler{}

type handlers struct {
	list []*Handler
}

var Handlers handlers

func (hs *handlers) Register(h *Handler) *handlers {
	if h == nil {
		panic("register: handlers should not be nil")
	}

	hs.list = append(hs.list, h)
	return hs
}

//func alias(node *node) string {
//	switch node.typ {
//	case "meta":
//		return node.typ
//	case "block":
//		return strings.Split(node.val, " ")[0]
//	default:
//		return "_" + node.typ
//	}
//}

func (hs *handlers) Match(node *node) (h *Handler) {
	for _, h := range hs.list {
		if h.Match(node) {
			return h
		}
	}
	return omitHandler
}
