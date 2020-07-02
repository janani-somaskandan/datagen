package utils

import(
	"io/ioutil"
	"os"
	"sync"
	
)

var file *os.File 
var m sync.Mutex


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

func (f FileWriter) RegisterOutputFile(FileName string){

	workingDirectory, _:= os.Getwd()
	path := workingDirectory +"/"+ FileName

	var _, err = os.Stat(path)

    // create file if not exists
    if os.IsNotExist(err) {
        var _, err = os.Create(path)
        if err != nil {
            return
        }
    }

	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if(err != nil){

	}
}