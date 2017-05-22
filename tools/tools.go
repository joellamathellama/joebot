package tools

import (
	"fmt"
	// "io"
	// "bytes"
	// "encoding/gob"
	// "bufio"
	// "reflect"
	"runtime"
	// "io/ioutil"
	"log"
	"os"
	"regexp"
	// "time"
)

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func RegexpMatch(regex string, word string) bool {
	res, err := regexp.MatchString(regex, word)
	if err != nil {
		WriteErr(err)
		fmt.Println(err)
	}
	return res
}

func MarkdownWrap(msg string) (mm string) {
	mm = fmt.Sprintf("```md\n%s\n```\n", msg)
	return
}

/*
	DEBUG LOGGING
*/

var path = "debug.txt"

func WriteLog(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if ok != true {
		file = "Unknown File"
		line = 0
	}
	logMsg := fmt.Sprintf("%s %d: %s", file, line, msg)

	if !FileExists(path) {
		CreateFile(path)
	}

	// f, err := os.Create(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		checkError(err)
	}
	defer f.Close()

	// w := bufio.NewWriter(f)
	// _, err = w.WriteString(logMsg)
	// checkError(err)

	// w.Flush()

	log.SetFlags(log.Ldate | log.Ltime) // Llongfile/Lshortfile only showed tools.go #, not where it was called
	log.SetOutput(f)
	log.Println(logMsg)
}

func WriteErr(err error) {
	errString := err.Error()
	_, file, line, ok := runtime.Caller(1)
	if ok != true {
		file = "Unknown File"
		line = 0
	}
	errMsg := fmt.Sprintf("%s %d: %s", file, line, errString)

	if !FileExists(path) {
		CreateFile(path)
	}

	// f, err := os.Create(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		checkError(err)
	}
	defer f.Close()

	// w := bufio.NewWriter(f)
	// _, err = w.WriteString(errMsg)
	// checkError(err)

	// w.Flush()

	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(f)
	log.Println(errMsg)
}

// FileExists checks for a file's existence
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreateFile creates a file
func CreateFile(name string) error {
	fo, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		fo.Close()
	}()
	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(666)
	}
}
