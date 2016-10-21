package main

import (
	"regexp"
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
		panic(err)
	} else {
		return res
	}
}
