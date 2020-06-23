package config

import(
	"./parser"
	"../utils"
	"time"
	"fmt"
)

var Config Configuration
func GenerateInputConfig(parserInstance parser.IParser, FileName string){

	InputConfig := parserInstance.Parse(utils.ReadFile(FileName), Config)
	Config = *InputConfig.(*Configuration)
	for _, element := range Config.User_segments {
		if element.Start_Time.IsZero() {
			element.Start_Time = time.Now().UTC()
		}
	}
	fmt.Println(Config)
}