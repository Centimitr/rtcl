package main

func lexArticleBody(l *lexer) stateFn {
	l.trim()
	return lexBlock(l)
}

func lexBlock(l *lexer) stateFn {
	if l.peek() == eof {
		return nil
	}

	if l.startWithBlankLine() {
		l.emit(itemBlankLine)
	}
	l.trim()

	// if block end
	if l.startWith("}") {
		l.emit(itemBlockRight)
		l.ignoreNext()
		l.ignoreLineEnd()
		return lexBlock(l)
	}
	//	current line if cmd or text
	if l.untilMatchOrLineEnd('{') {
		l.emit(itemCmd)
		l.ignoreNext()
		l.ignoreLineEnd()
		l.emit(itemBlockLeft)
		return lexBlock(l)
	}
	// text line if empty
	if l.active() != "" {
		l.emit(itemText)
		l.ignoreLineEnd()
		l.emit(itemSep)
	}
	return lexBlock(l)
}
