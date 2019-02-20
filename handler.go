package rtcl

type HandleChildrenFn func()
type HandleFn func(node *node, handleChildren HandleChildrenFn)

type Handler struct {
	Match  func(node *node) bool
	Handle HandleFn
}

var omitHandler = &Handler{}

type handlers struct {
	list           []*Handler
	blockHandleFns map[string]HandleFn
}

var DefaultHandlers handlers

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

func (hs *handlers) RegisterBlockType(typ string, fn HandleFn) *handlers {
	if hs.blockHandleFns == nil {
		hs.blockHandleFns = make(map[string]HandleFn)
	}
	hs.blockHandleFns[typ] = fn
	return hs
}

func (hs *handlers) Match(node *node) (h *Handler) {
	for _, h := range hs.list {
		if h.Match(node) {
			return h
		}
	}
	return omitHandler
}

func init() {
	DefaultHandlers.
		RegisterType("block.command", func(node *node, handleChildren HandleChildrenFn) {
			node.representation = "WRAPPER TYPE: " + node.val
		}).
		RegisterType("block", func(node *node, handleChildren HandleChildrenFn) {
			if node.child == nil || node.child.typ != "block.command" {
				return
			}
			args := NewArgsFromString(node.child.val)
			fn := DefaultHandlers.blockHandleFns[args.First]

			if fn != nil {
				fn(node, handleChildren)
			}
		})
}
