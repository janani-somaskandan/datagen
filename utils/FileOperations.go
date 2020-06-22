package utils

import(
	"io/ioutil"
	"os"
	"sync"
)

var file *os.File 
var m sync.Mutex

func RegisterOutputFile(FileName string){

	workingDirectory, _:= os.Getwd()
	var err error
	file, err = os.OpenFile(workingDirectory +"/"+ FileName, os.O_APPEND|os.O_WRONLY, 0644)
	if(err != nil){

	}
}

func ReadFile(FileName string)[]byte {

	workingDirectory, err := os.Getwd()
	data,err := ioutil.ReadFile(workingDirectory + "/config/" + FileName)
	if(err != nil){
		
	}
	return data

}

type FileWriter struct{}
func (f FileWriter) Write(data string){
	m.Lock()
		file.WriteString(data + "\n")
	m.Unlock()
}