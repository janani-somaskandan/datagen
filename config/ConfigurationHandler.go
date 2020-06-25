package config

import(
	"./parser"
	"../utils"
	"time"
)

var ConfigV2 ConfigurationV2
func GenerateInputConfigV2(parserInstance parser.IParser, FileName string){

	InputConfig := parserInstance.Parse(utils.ReadFile(FileName), ConfigV2)
	ConfigV2 = *InputConfig.(*ConfigurationV2)
	for _, element := range ConfigV2.User_segments {
		if element.Start_Time.IsZero() {
			element.Start_Time = time.Now().UTC()
		}
	}
}