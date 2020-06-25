package registration

import(
	"fmt"
	"../config"
	"../config/parser"
	"../utils"
	"../adaptors"
)

var WriterInstance adaptors.Writer

func RegisterHandlers(){	
	fmt.Println("Registering Handlers")
	var _parser parser.IParser
	_parser = parser.YamlParser{}
	config.GenerateInputConfigV2(_parser,"sampleconfigV2.yaml")
	WriterInstance = utils.FileWriter{}
    utils.RegisterOutputFile("output.txt")
}