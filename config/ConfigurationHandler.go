package config

import(
	"./parser"
	"../utils"
	"time"
	"fmt"
)

var ConfigV2 ConfigurationV2
func GenerateInputConfigV2(parserInstance parser.IParser, FileName string){

	InputConfig := parserInstance.Parse(utils.ReadFile(FileName), ConfigV2)
	ConfigV2 = *InputConfig.(*ConfigurationV2)
	for item, element := range ConfigV2.User_segments {
		if element.Start_Time.IsZero() {
			element.Start_Time = time.Now().UTC()
			ConfigV2.User_segments[item] = element
		}
	}
	fmt.Println(ConfigV2)
}