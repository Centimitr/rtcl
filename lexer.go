package main

import (
	"strings"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	state stateFn
	pos   int
	start int
	width int
	items chan item
}

func newLexer(input string) *lexer {
	return &lexer{
		input: input,
		state: lexArticle,
		items: make(chan item),
	}
}

func (l *lexer) run() {
	//l.input += "\n"
	for l.state != nil {
		l.state = l.state(l)
	}
	l.items <- item{typ: itemEOF}
	close(l.items)
}

func (l *lexer) debug(s string) {
	l.items <- item{typ: itemError, val: s}
}

// emit

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.active()}
	l.start = l.pos
}

// content

func (l *lexer) active() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) rest() string {
	return l.input[l.start:]
}

func (l *lexer) startWith(s string) bool {
	return strings.HasPrefix(l.rest(), s)
}

func (l *lexer) startWithLine(s string) bool {
	return l.startWith(s + "\n")
}

func (l *lexer) startWithBlankLine() bool {
	for _, r := range l.rest() {
		switch r {
		case ' ':
		case '\n':
			return true
		default:
			return false
		}
	}
	return false
}
