package registration

/*
Registering config Reader, ouput writer
*/

import(
	"../config"
	"../config/parser"
	"../utils"
	"../adaptors"
	Log "../utils/Log"
)

var WriterInstance adaptors.Writer

func RegisterHandlers(){	

	Log.Debug.Println("Registering Handlers")

	Log.Debug.Println("Registering Yaml Parser")
	var _parser parser.IParser
	_parser = parser.YamlParser{}
	config.GenerateInputConfigV2(_parser,"/../config/livspace.yaml")
	// log.Println("Registering Output to File Writer")
	// WriterInstance = utils.FileWriter{}
	// WriterInstance.RegisterOutputFile(config.ConfigV2.Output_file_name)
	Log.Debug.Println("Registering Output to Log Writer")
	WriterInstance = utils.LogWriter{}
	WriterInstance.RegisterOutputFile(config.ConfigV2.Output_file_name)

	Log.Debug.Println("Registering UserData to Log Writer")
	WriterInstance.RegisterUserDataFile(config.ConfigV2.User_data_file_name)
	
	Log.Debug.Println("Registration Done !!!")
}