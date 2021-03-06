package utils

/*
Util for File based operations
*/

import(
	"io/ioutil"
	"os"
	"sync"
	Log "../utils/Log"
	"compress/gzip"
	"path/filepath"
	"strings"
)

var file *os.File 
var m sync.Mutex


func GetAllUnreadFiles(FileRootPath string, filenameprefix string)[]string {

	var files []string
	Log.Debug.Printf("Searching for all the files in path %s with extension .gz and File name prefix %s", FileRootPath,filenameprefix)

    err := filepath.Walk(FileRootPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".gz" {
			return nil
		}
		if(!strings.HasPrefix(info.Name(), filenameprefix)){
			return nil
		}
        files = append(files, path)
        return nil
    })
    if err != nil {
        Log.Error.Fatal(err)
    }
    return files
}

func GetFileHandle(FilePath string) (*os.File){

	Log.Debug.Printf("ReadingFile %s", FilePath)
	handle, err := os.Open(FilePath)
	if err != nil {
		Log.Error.Fatal(err)
	}
	return handle
}

func GetFileHandlegz(FilePath string) (*gzip.Reader){

	Log.Debug.Printf("ReadingFile %s", FilePath)
	handle, err := os.Open(FilePath)
	if err != nil {
		Log.Error.Fatal(err)
	}
	zipReader, err := gzip.NewReader(handle)
	if err != nil {
		Log.Error.Fatal(err)
	}
	defer zipReader.Close()
	return zipReader
}

func ReadFile(FileName string)[]byte {

	workingDirectory, err := os.Getwd()
	data,err := ioutil.ReadFile(workingDirectory + FileName)
	if(err != nil){
		Log.Error.Fatal(err)
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
		Log.Debug.Printf("Creating File %s", path)
        var _, err = os.Create(path)
        if err != nil {
            Log.Error.Fatal(err)
        }
    }

	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if(err != nil){
		Log.Error.Fatal(err)
	}
}

func CreateDirectoryIfNotExists(folderPath string){
	_, err := os.Stat(folderPath)
 
	if os.IsNotExist(err) {
		dir := os.MkdirAll(folderPath, 0755)
		if dir != nil {
			Log.Error.Fatal(err)
		}
	}
}

func MoveFiles(old string, new string){
	err := os.Rename(old, new)
	if err != nil {
		Log.Error.Fatal(err)
	}
}