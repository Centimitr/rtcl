package rtcl

func lexArticleBody(l *lexer) stateFn {
	l.trim()
	return lexBlock(l)
}

func handleBlockStart(l *lexer, ignoreBlankLine bool) (stateFn, bool) {
	if l.peek() == eof {
		return nil, true
	}

	if l.startWithBlankLine() {
		if !ignoreBlankLine {
			l.emit(itemBlankLine)
		}
	}
	l.trim()

	// if block end
	if l.startWith("}") {
		l.emit(itemBlockRight)
		l.ignoreNext()
		l.ignoreLineEnd()
		return lexBlock(l), true
	}
	return nil, false
}

func lexBlock(l *lexer) stateFn {
	if fn, needRtn := handleBlockStart(l, false); needRtn {
		return fn
	}

	//	current line if cmd or text
	if l.untilMatchOrLineEnd('{') {
		metaBlock := l.trimmedActive() == "meta"
		l.emitWithTrim(itemCmd)
		l.ignoreNext()
		l.ignoreLineEnd()
		l.emit(itemBlockLeft)
		if metaBlock {
			return lexMetaItem(l)
		}
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

func lexMetaItem(l *lexer) stateFn {
	if fn, needRtn := handleBlockStart(l, true); needRtn {
		return fn
	}
	l.trim()
	l.untilMatchOrLineEnd(',')
	l.emitWithTrim(itemMetaItem)
	l.ignoreNext()
	return lexMetaItem(l)
}
