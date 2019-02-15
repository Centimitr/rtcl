package main

type itemType string

type item struct {
	typ itemType
	val string
}

const (
	itemError itemType = "itemError"
	itemEOF            = "itemEOF"
	itemBack           = "itemBack"
	itemBlock          = "itemBlock"

	itemBlockArticle = "itemBlockArticle"

	itemBlockMeta     = "itemBlockMeta"
	itemBlockMetaArgs = "itemBlockMetaArgs"
	itemBlockMetaKVs  = "itemBlockMetaKVs"
	itemMetaArg       = "itemMetaArg"
	itemMetaSep       = "itemMetaSep"
	itemMetaKV        = "itemMetaKV"
	itemMetaKey       = "itemMetaKey"
	itemMetaValue     = "itemMetaValue"

	//itemBlockBody = "itemBlockBody"
	itemBlockLeft      = "itemBlockLeft"
	itemBlockRight     = "itemBlockRight"
	itemText           = "itemText"
	itemSep            = "itemSep"
	itemBlankLine      = "itemBlankLine"
	itemCmd            = "itemCmd"
	itemDecoratorLeft  = "itemDecoratorLeft"
	itemDecoratorRight = "itemDecoratorRight"
)
