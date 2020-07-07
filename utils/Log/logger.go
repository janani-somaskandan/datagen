package Log

/*
Log registration for debug and error logs
*/

import(
	"os"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

var Debug *log.Logger
var Error *log.Logger

func RegisterLogFiles(){

    fmt.Println("Check For all debug logs in: debugsLogs.log")
	debuglog, err := os.OpenFile("debugLogs.log",  os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("error opening file: %v", err)
        os.Exit(1)
    }

	Debug = log.New(debuglog, "", log.LstdFlags)
    Debug.SetOutput(&lumberjack.Logger{
		Filename:   "debugLogs.log",
		MaxSize:    1, // megabytes
		MaxAge:     10, // days
		Compress:   true, // disabled by default
    })
    
    fmt.Println("Check For all error logs in: errorLogs.log")
	errorlog, err := os.OpenFile("errorLogs.log",  os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("error opening file: %v", err)
        os.Exit(1)
    }

    Error = log.New(errorlog, "", log.LstdFlags)
    Error.SetOutput(&lumberjack.Logger{
		Filename:   "errorLogs.log",
		MaxSize:    1, // megabytes
		MaxAge:     10, // days
		Compress:   true, // disabled by default
    })   
}