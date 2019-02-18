package main

import "fmt"

func init() {
	Handlers.
		RegisterInternal("blank", nil).
		RegisterInternal("text", nil).
		RegisterInternal("sep", nil).

		Register("meta", &Handler{
			Pre: func(node *node) bool {
				fmt.Println("META:", node.typ, astChildren(node))
				return false
			}}).

		Register("#", nil).
		Register("-", nil).
		Register("[]", nil).
		Register("define", nil)
}
