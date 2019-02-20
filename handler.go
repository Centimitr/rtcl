package rtcl

type HandleChildrenFn func()
type HandleFn func(node *node, handleChildren HandleChildrenFn)

type Handler struct {
	Match  func(node *node) bool
	Handle HandleFn
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

func NewTypeHandler(typ string, fn HandleFn) *Handler {
	if fn == nil {
		fn = func(node *node, handleChildren HandleChildrenFn) {}
	}
	return &Handler{
		Match: func(node *node) bool {
			return node.typ == typ
		},
		Handle: fn,
	}
}

func (hs *handlers) RegisterType(typ string, fn HandleFn) *handlers {
	return hs.Register(NewTypeHandler(typ, fn))
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
