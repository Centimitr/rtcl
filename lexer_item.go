package rtcl

type itemType string

type item struct {
	typ itemType
	val string
}

const (
	itemError itemType = "itemError"
	itemEOF            = "itemEOF"

	itemMetaArg   = "itemMetaArg"
	itemMetaSep   = "itemMetaSep"
	itemMetaKey   = "itemMetaKey"
	itemMetaValue = "itemMetaValue"
	itemMetaItem  = "itemMetaItem"

	itemRaw        = "itemRaw"
	itemBlockLeft  = "itemBlockLeft"
	itemBlockRight = "itemBlockRight"
	itemText       = "itemText"
	//itemSep            = "itemSep"
	itemBlankLine      = "itemBlankLine"
	itemCmd            = "itemCmd"
	itemDecoratorLeft  = "itemDecoratorLeft"
	itemDecoratorRight = "itemDecoratorRight"
)
