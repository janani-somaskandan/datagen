package parser

type IParser interface{
	Parse([]byte, interface{})(interface{})
}