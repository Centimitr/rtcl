package rtcl

import (
	"unicode/utf8"
)

// move pos

const eof = rune(-1)

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.rest())
	l.pos += l.width
	return r
}

func (l *lexer) peek() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *lexer) untilMatch(r rune) {
	for {
		switch l.peek() {
		case r, eof:
			return
		default:
			l.next()
		}
	}
}

func (l *lexer) untilLineEnd() {
	l.untilMatch('\n')
}

func (l *lexer) untilMatchOrLineEnd(r rune) bool {
	for {
		switch l.peek() {
		case r:
			return true
		case '\n', eof:
			return false
		default:
			l.next()
		}
	}
}

func (l *lexer) untilSeeString(s string) {
	for {
		if l.waitingStartWith(s) {
			println("OK")
			break
		}
		if l.next() == eof {
			break
		}
	}
}

func (l *lexer) untilMatchLine(s string) {
	for {
		if l.startWith(s + "\n") {
			break
		}
		if l.next() == eof {
			break
		}
	}
}

func (l *lexer) trim() {
	for {
		switch l.peek() {
		case ' ', '	', '\n':
			l.next()
		case eof:
			return
		default:
			l.start = l.pos
			return
		}
	}
}

func (l *lexer) ignoreCurrent() {
	l.start = l.pos
}

func (l *lexer) ignoreNext() {
	l.next()
	l.ignoreCurrent()
}

func (l *lexer) ignoreWhitespace() {
	if l.peek() == ' ' || l.peek() == '	' {
		l.next()
		l.ignoreCurrent()
	}
}

func (l *lexer) ignoreLineEnd() {
	if l.peek() == '\n' {
		l.next()
		l.ignoreCurrent()
	}
}

func (l *lexer) ignoreLine() {
	l.untilLineEnd()
	l.next()
	l.ignoreCurrent()
}
