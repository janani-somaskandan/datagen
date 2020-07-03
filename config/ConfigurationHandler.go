package config

/*
Reads the input config and applies required customization on object
*/

import(
	"./parser"
	"../utils"
	"time"
	Log "../utils/Log"
)

var ConfigV2 ConfigurationV2
func GenerateInputConfigV2(parserInstance parser.IParser, FileName string){

	Log.Debug.Printf("Processing input config - %s", FileName)
	InputConfig := parserInstance.Parse(utils.ReadFile(FileName), ConfigV2)
	ConfigV2 = *InputConfig.(*ConfigurationV2)
	for item, element := range ConfigV2.User_segments {
		if element.Start_Time.IsZero() {
			element.Start_Time = time.Now().UTC()
			ConfigV2.User_segments[item] = element
		}
	}
	Log.Debug.Printf("Operating with config %v", ConfigV2)
}