package main

func lexArticle(l *lexer) stateFn {
	return lexArticleMeta(l)
}

func lexArticleMeta(l *lexer) stateFn {
	return lexMetaArg(l)
}

func lexMetaArg(l *lexer) stateFn {
	if l.peek() == eof {
		return nil
	}

	l.trim()
	if l.startWithLine("===") {
		l.emit(itemMetaSep)
		l.ignoreLine()
		return lexMetaKV(l)
	}
	l.untilLineEnd()
	l.emit(itemMetaArg)
	l.trim()
	return lexMetaArg(l)
}

func lexMetaKV(l *lexer) stateFn {
	if l.peek() == eof {
		return nil
	}

	if l.startWithBlankLine() {
		l.ignoreLine()
		l.emit(itemBlankLine)
		return lexArticleBody(l)
	}
	l.untilMatch(' ')
	l.emit(itemMetaKey)
	l.untilLineEnd()
	l.emit(itemMetaValue)
	l.ignoreLineEnd()
	return lexMetaKV(l)
}
