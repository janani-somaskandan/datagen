package main

/*
main file of datagen job
*/

import(
	"../registration"
	"../operations"	
	Log "../utils/Log"
)


func main(){
	Log.RegisterLogFiles()
	registration.RegisterHandlers()
	operations.OperateV1()
}  
