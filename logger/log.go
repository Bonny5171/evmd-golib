package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	InfoLog    *log.Logger
	TraceLog   *log.Logger
	DebugLog   *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
	FatalLog   *log.Logger
	PanicLog   *log.Logger
)

func Init(appname string, infoHandle, traceHandle, debugHandle, warningHandle, errorHandle io.Writer) {
	if len(appname) > 0 {
		appname = fmt.Sprintf("[%s] ", appname)
	}

	InfoLog = log.New(infoHandle, fmt.Sprintf("%sINFO   : ", appname), log.Ldate|log.Ltime)
	TraceLog = log.New(traceHandle, fmt.Sprintf("%sTRACE  : ", appname), log.Ldate|log.Ltime)
	DebugLog = log.New(debugHandle, fmt.Sprintf("%sDEBUG  : ", appname), log.Ldate|log.Ltime)
	WarningLog = log.New(warningHandle, fmt.Sprintf("%sWARNING: ", appname), log.Ldate|log.Ltime)
	ErrorLog = log.New(errorHandle, fmt.Sprintf("%sERROR  : ", appname), log.Ldate|log.Ltime)
	FatalLog = log.New(errorHandle, fmt.Sprintf("%sFATAL  : ", appname), log.Ldate|log.Ltime)
	PanicLog = log.New(errorHandle, fmt.Sprintf("%sPANIC  : ", appname), log.Ldate|log.Ltime)
}

func debugActive() bool {
	da, err := strconv.ParseBool(os.Getenv("DEBUGGO"))
	if err != nil {
		return false
	}
	return da
}

func Debug(t interface{}) {
	if debugActive() {
		DebugLog.Print(t)
	}
}

func Debugf(format string, v ...interface{}) {
	if debugActive() {
		DebugLog.Printf(format, v...)
	}
}

func Debugln(t ...interface{}) {
	if debugActive() {
		DebugLog.Println(t...)
	}
}

func traceActive() bool {
	ta, err := strconv.ParseBool(os.Getenv("TRACEGO"))
	if err != nil {
		return true
	}
	return ta
}

func Trace(t interface{}) {
	if traceActive() {
		TraceLog.Print(t)
	}
}

func Tracef(format string, v ...interface{}) {
	if traceActive() {
		TraceLog.Printf(format, v...)
	}
}

func Traceln(t ...interface{}) {
	if traceActive() {
		TraceLog.Println(t...)
	}
}

func Info(i interface{}) {
	InfoLog.Print(i)
}

func Infof(format string, v ...interface{}) {
	InfoLog.Printf(format, v...)
}

func Infoln(i ...interface{}) {
	InfoLog.Println(i...)
}

func Warning(w interface{}) {
	WarningLog.Print(w)
}

func Warningf(format string, v ...interface{}) {
	WarningLog.Printf(format, v...)
}

func Warningln(w ...interface{}) {
	WarningLog.Println(w...)
}

func Error(e interface{}) {
	ErrorLog.Print(e)
}

func Errorf(format string, v ...interface{}) {
	ErrorLog.Printf(format, v...)
}

func Errorln(e ...interface{}) {
	ErrorLog.Println(e...)
}

func Fatal(e interface{}) {
	FatalLog.Fatal(e)
}

func Fatalf(format string, v ...interface{}) {
	FatalLog.Fatalf(format, v...)
}

func Fatalln(e ...interface{}) {
	FatalLog.Fatalln(e...)
}

func Panic(e interface{}) {
	PanicLog.Panic(e)
}

func Panicf(format string, v ...interface{}) {
	PanicLog.Panicf(format, v...)
}

func Panicln(e ...interface{}) {
	PanicLog.Panicln(e...)
}
