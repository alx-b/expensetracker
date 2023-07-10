package logger

// Package logger implements a barebone logger writing to logs.txt file.
//
// Import and use its functions with specific log level.
// Defer CloseFile() in your main *(not required but better safe than sorry!)

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var loggers Loggers = createLoggers()

// Struct to gather all pointers to log.Logger and os.File
type Loggers struct {
	file  *os.File
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	panic *log.Logger
}

// createLoggers returns a Loggers struct
func createLoggers() Loggers {
	file := openFile()

	info := log.New(file, "INFO: ", log.LstdFlags)
	warn := log.New(file, "WARN: ", log.LstdFlags)
	error := log.New(file, "ERROR: ", log.LstdFlags)
	panic := log.New(file, "PANIC: ", log.LstdFlags)

	return Loggers{
		info:  info,
		warn:  warn,
		error: error,
		panic: panic,
		file:  file,
	}
}

// openFile opens and returns a pointer to a file.
func openFile() *os.File {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

// CloseFile syncs and closes file.
func CloseFile() {
	if err := loggers.file.Sync(); err != nil {
		Panic("file.Sync failing")
	}

	if err := loggers.file.Close(); err != nil {
		Panic("file.Close failing")
	}
}

// printLineToFile takes in a text string and print function
// uses runtime.Caller(2) for accessing the origin function caller's properties
// and write to file with print callback function.
func printLineToFile(text string, print func(string)) {
	programCounter, fullPath, lineNumber, ok := runtime.Caller(2)
	if !ok {
		Panic("runtime.Caller is not ok")
		return
	}

	fullFunctionPath := runtime.FuncForPC(programCounter).Name()
	functionName := getLastElementOfStringPath(fullFunctionPath)
	fileName := getLastElementOfStringPath(fullPath)

	line := fmt.Sprintf("%s:%d:%s %s", fileName, lineNumber, functionName, text)
	print(line)
}

// getLastElementOfStringPath returns lastElement of a path string.
func getLastElementOfStringPath(path string) string {
	//TODO Might need to have a window version with "\"?
	splittedString := strings.Split(path, "/")
	lastElement := splittedString[len(splittedString)-1]
	return lastElement
}

// Info takes in a string and passes in a log.Logger.Println function to
// printLineToFile function.
func Info(text string) {
	printLineToFile(text, func(t string) { loggers.info.Println(t) })
}

// Warn takes in a string and passes in a log.Logger.Println function to
// printLineToFile function.
func Warn(text string) {
	printLineToFile(text, func(t string) { loggers.warn.Println(t) })
}

// Error takes in a string and passes in a log.Logger.Println function to
// printLineToFile function.
func Error(text string) {
	printLineToFile(text, func(t string) { loggers.error.Println(t) })
}

// Panic takes in a string and passes in a log.Logger.Println function to
// printLineToFile function.
func Panic(text string) {
	printLineToFile(text, func(t string) { loggers.panic.Println(t) })
}
