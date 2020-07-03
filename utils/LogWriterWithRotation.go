package utils

/* 
Util for Log operations with Log Rotation
*/

import(
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"fmt"
	"os"
)

var fLogger *log.Logger
type LogWriter struct{}

func (f LogWriter) RegisterOutputFile(FileName string){

	File, err := os.OpenFile(FileName,  os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("error opening file: %v", err)
        os.Exit(1)
    }

	fLogger = log.New(File, "", log.LstdFlags)
    fLogger.SetOutput(&lumberjack.Logger{
		Filename:   FileName,
		MaxSize:    1, // megabytes
		MaxAge:     10, // days
		Compress:   true, // disabled by default
    })
	// SET MaxBackups: 2 if required
}

func (f LogWriter) Write(data string){
	fLogger.Printf("%s\n",data)
}