package main

import (
	"fmt"
	// "io"
	// "bytes"
	// "encoding/gob"
	// "bufio"
	// "io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
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

func writeLog(logString string) {
	if !FileExists(path) {
		CreateFile(path)
	}

	today := time.Now()
	finalLog := fmt.Sprintf("%s: %s", today, logString)

	// f, err := os.Create(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		checkError(err)
	}
	defer f.Close()

	// w := bufio.NewWriter(f)
	// _, err = w.WriteString(logString)
	// checkError(err)

	// w.Flush()

	log.SetOutput(f)
	log.Println(finalLog)
}

func writeErr(err error) {
	logString := err.Error()
	today := time.Now()
	logString = fmt.Sprintf("%s: %s", today, logString)

	f, err := os.Create(path)
	if err != nil {
		checkError(err)
	}
	defer f.Close()

	// w := bufio.NewWriter(f)
	// _, err = w.WriteString(logString)
	// checkError(err)

	// w.Flush()

	log.SetOutput(f)
	log.Println(logString)
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
