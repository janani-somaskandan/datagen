package parser

import (
	"gopkg.in/yaml.v2"
	"reflect"
)

type YamlParser struct{}
func (y YamlParser) Parse(FileContents []byte, outputObj interface{}) (interface{}) {

    obj := reflect.New(reflect.TypeOf(outputObj)).Interface()
	err := yaml.Unmarshal(FileContents, obj)
	if(err != nil){
		panic(err)
	}
	return obj
}