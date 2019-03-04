package rtcl

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

func emit(l *lexer, t itemType, active string) {
	l.items <- item{t, active}
	//fmt.Println("emit:", t, active, len(active))
	l.start = l.pos
}

func (l *lexer) emit(t itemType) {
	emit(l, t, l.active())
}

func (l *lexer) emitWithTrim(t itemType) {
	emit(l, t, l.trimmedActive())
}

// content

func (l *lexer) active() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) trimmedActive() string {
	start := l.start
	pos := l.pos
	for {
		switch l.input[start] {
		case ' ', '	':
			start++
			continue
		}
		break
	}
	for {
		switch l.input[pos-1] {
		case ' ', '	':
			pos--
			continue
		}
		break
	}
	return l.input[start:pos]
}

func (l *lexer) rest() string {
	return l.input[l.start:]
}

func (l *lexer) waiting() string {
	return l.input[l.pos:]
}

func (l *lexer) startWith(s string) bool {
	return strings.HasPrefix(l.rest(), s)
}

func (l *lexer) waitingStartWith(s string) bool {
	return strings.HasPrefix(l.waiting(), s)
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
