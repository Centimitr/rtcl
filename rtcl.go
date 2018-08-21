package rtcl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
)

// lex

type itemType int

const (
	itemError itemType = iota
	itemCommand
	itemLineEnd
	itemBlockLeft
	itemBlockRight
	itemInlineBlockLeft
	itemInlineBlockRight
	itemEOF
)

type item struct {
	typ itemType
	val string
}

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
		state: lexCommand,
		items: make(chan item, 2),
	}
}

func (l *lexer) run() {
	l.input += "\n"
	for l.state != nil {
		l.state = l.state(l)
	}
	close(l.items)
}

const eof = rune(-1)

func (l *lexer) emit(t itemType) {
	l.items <- item{typ: t, val: strings.TrimSpace(l.active())}
	l.start = l.pos
}

func (l *lexer) active() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func lexCommand(l *lexer) stateFn {
	cnt := 0
	for {
		r := l.next()
		cnt++
		if r == eof {
			break
		}
		switch r {
		case '\n':
			l.emit(itemCommand)
			l.emit(itemLineEnd)
			return lexCommand
		case '{':
			l.backup()
			l.emit(itemCommand)
			l.next()
			return lexBlockStart
		case '}':
			l.backup()
			l.emit(itemCommand)
			l.next()
			return lexBlockEnd
		}
	}
	if l.pos > l.start {
	}
	l.emit(itemEOF)
	return nil
}

func lexBlockStart(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case '{':
			l.emit(itemInlineBlockLeft)
			return lexCommand
		case '\n':
			l.emit(itemBlockLeft)
			l.emit(itemLineEnd)
			return lexCommand
		}
	}
}

func lexBlockEnd(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case '}':
			l.emit(itemInlineBlockRight)
			return lexCommand
		case '\n':
			l.emit(itemBlockRight)
			l.emit(itemLineEnd)
			return lexCommand
		}
	}
}

// parse

type transform func()

type transforms map[string]transform

func funcName(fn interface{}) string {
	p := reflect.ValueOf(fn).Pointer()
	path := runtime.FuncForPC(p).Name()
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

func Register(name string, namespace string, fns ...transform) {
	var _, _ transforms
	switch name {
	case SyntaxNameMeta:
	case SyntaxNameContent:
	default:
		panic("name undefined")
	}
}

type syntax struct {
	name     string
	value    string
	parent   *syntax
	children []*syntax
}

func (s *syntax) child(name string) *syntax {
	for _, syntax := range s.children {
		if syntax.name == name {
			return syntax
		}
	}
	return nil
}

func (s *syntax) lastChild() *syntax {
	if l := len(s.children); l >= 1 {
		return s.children[l-1]
	}
	return nil
}

func (s *syntax) addChild(syn *syntax) *syntax {
	syn.parent = s
	s.children = append(s.children, syn)
	return syn
}

func (s *syntax) fill(item item, t transforms) *syntax {
	return nil
}

func (s *syntax) print(indent int) {
	spaces := strings.Repeat("\t", indent)
	name := s.name
	if name == "" {
		name = "NIL"
	}
	fmt.Printf("%s [%s] %s\n", spaces, name, s.value)
	for _, child := range s.children {
		child.print(indent + 1)
	}
}

const (
	syntaxNameRoot    = "ROOT"
	SyntaxNameMeta    = "META"
	SyntaxNameContent = "CONTENT"
)

func newSyntaxTree() *syntax {
	return &syntax{
		name:   syntaxNameRoot,
		parent: nil,
		children: []*syntax{
			{name: SyntaxNameMeta},
			{name: SyntaxNameContent},
		},
	}
}

type inlineItems []item

func (it *inlineItems) needOmitLine() bool {
	return len(*it) == 1 && (*it)[0].val == ""
}

func (it *inlineItems) flush(node *syntax) {
	if node != nil {
		for _, itm := range *it {
			if itm.typ == itemCommand && itm.val == "" {
				continue
			}
			node.children = append(node.children, &syntax{parent: node, value: itm.val})
		}
	}
	*it = (*it)[:0]
}

func parse(items <-chan item) *syntax {
	t := newSyntaxTree()
	node := t.child(SyntaxNameMeta)
	var curLineItems inlineItems
	hasSwitchToContent := false
	for item := range items {
		switch item.typ {
		case itemCommand:
			curLineItems = append(curLineItems, item)
		case itemLineEnd:
			switch {
			case curLineItems.needOmitLine():
				curLineItems.flush(nil)
				if !hasSwitchToContent {
					node = t.child(SyntaxNameContent)
					hasSwitchToContent = true
				}
			default:
				curLineItems.flush(node)
			}
			if node.name == "LINE" {
				node = node.parent
			}
		case itemBlockLeft:
			curLineItems.flush(node)
			if len(node.children) < 1 {
				node.addChild(&syntax{name: "BLOCK"})
			}
			node = node.lastChild()
		case itemBlockRight:
			node = node.parent
		case itemInlineBlockLeft:
			container := node.addChild(&syntax{name: "LINE"})
			curLineItems.flush(container)
			node = container.addChild(&syntax{name: "INLINEBLOCK"})
		case itemInlineBlockRight:
			curLineItems.flush(node)
			node = node.parent
		}
	}
	return t
}

func Parse(input string) {
	l := newLexer(input)
	go l.run()
	t := parse(l.items)
	t.print(0)
}

func ParseFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	Parse(string(content))
	return nil
}
