package main

import "fmt"

type Handler struct {
	Pre  func(*node) bool
	Post func(*node)
}

var defaultHandler = &Handler{Pre: func(node *node) bool {
	if node.typ != "block" {
		fmt.Println("unhandle:", node.typ)
	}
	return true
}}
var omitHandler = &Handler{}

type handlers map[string]*Handler

var Handlers handlers

func (handlers *handlers) Register(typ string, h *Handler) *handlers {
	if *handlers == nil {
		*handlers = make(map[string]*Handler)
	}

	m := *handlers
	m[typ] = omitHandler
	if h != nil {
		m[typ] = h
	}
	return handlers
}

func (handlers *handlers) RegisterInternal(typ string, h *Handler) *handlers {
	return handlers.Register("_"+typ, h)
}

func (handlers *handlers) RegisterMeta(h *Handler) *handlers {
	return handlers
}

func alias(node *node) string {
	switch node.typ {
	case "meta":
		return node.typ
	case "block":
		return node.val
	default:
		return "_" + node.typ
	}
}

func (handlers *handlers) Match(node *node) (h *Handler) {
	h, ok := (*handlers)[alias(node)]
	if !ok {
		h = defaultHandler
	}
	return h
}
