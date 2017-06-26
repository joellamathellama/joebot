package routes

import (
	"joebot/rds"
	"joebot/tools"
	"strings"
)

func myNotes(sender string, url string) (res string) {
	lowName := strings.ToLower(sender)
	res = ""
	if url == "GET" {
		// redis get persons team image
		link, err := rds.RedisGet(rds.RC, lowName)
		if err != nil {
			tools.WriteErr(err)
			res = "No note has been set!"
		} else {
			res = link
		}
	} else {
		// else set the url
		ok := rds.RedisSet(rds.RC, lowName, url)
		if !ok {
			tools.WriteLog("Error: myNotes() redisSet() Fail")
			res = "Something went wrong, alert the Master Llama!"
		} else {
			res = "Note set!"
		}
	}
	return
}

func getNotes(user string) (res string) {
	// redis get persons team image
	lowName := strings.ToLower(user)
	res, err := rds.RedisGet(rds.RC, lowName)
	if err != nil {
		tools.WriteErr(err)
		res = "That person has no note set!"
	}
	return
}
