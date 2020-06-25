package main

import(
	"./registration"
	"./operations"	
)


func main(){
	registration.RegisterHandlers()
	operations.OperateV2()
	
}  

