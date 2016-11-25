package main

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

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func regexpMatch(regex string, word string) bool {
	res, err := regexp.MatchString(regex, word)
	if err != nil {
		writeErr(err)
		fmt.Println(err)
	}
	return res
}

/*
	DEBUG LOGGING
*/

var path = "debug.txt"

func writeLog(msg string) {
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

func writeErr(err error) {
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

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

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
