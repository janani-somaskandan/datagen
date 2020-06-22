package config

import(
	"./parser"
	"../utils"
	"time"
)

var Config Configuration
func GenerateInputConfig(parserInstance parser.IParser, FileName string){

	InputConfig := parserInstance.Parse(utils.ReadFile(FileName), Config)
	Config = *InputConfig.(*Configuration)
	if Config.Start_Time.IsZero() {
		Config.Start_Time = time.Now().UTC()
	}
}