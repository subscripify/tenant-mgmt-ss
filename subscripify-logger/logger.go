package subscripifylogger

import (
	"log"
	"os"
)

var (
	TraceLog   *log.Logger
	DebugLog   *log.Logger
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
	FatalLog   *log.Logger
)

func init() {
	// file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	stdOut := os.Stdout

	TraceLog = log.New(stdOut, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLog = log.New(stdOut, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog = log.New(stdOut, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(stdOut, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(stdOut, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLog = log.New(stdOut, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}
