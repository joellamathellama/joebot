package routes

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/rds"
	"joebot/tools"
	"strings"
)

func storyRouteSS(s *dg.Session, playerName string) {
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "Story").Result()
	if err != nil {
		tools.WriteErr(err)
		messageSend(s, "Player's story not found! Try, idk, typing it correctly?")
	} else {
		messageSend(s, res)
	}
}

func slotesRouteSS(s *dg.Session, playerName string) {
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "Stones").Result()
	if err != nil {
		tools.WriteErr(err)
		messageSend(s, "Player's stones not found. Prolly cause you're stoned...")
	} else {
		messageSend(s, res)
	}
}

func ssherderRouteSS(s *dg.Session, playerName string) {
	// https://ssherder.com/characters/ID/
	// lookup player ID, add to URL, send message
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "ID").Result()
	if err != nil {
		tools.WriteErr(err)
		messageSend(s, "Who?!")
	} else {
		messageSend(s, "https://ssherder.com/characters/"+res)
	}
}

func skillsRouteSS(s *dg.Session, playerName string) {
	// https://ssherder.com/characters/ID/
	// lookup player ID, add to URL, send message
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "Skills").Result()
	if err != nil {
		tools.WriteErr(err)
		messageSend(s, "Player's skills not found. Sharpen your typing skills first...")
	} else {
		messageSend(s, res)
	}
}

func stoneRouteSS(s *dg.Session, stoneName string) {
	// rds.RedisGet(key)
	// if else error
	stoneTitle := strings.Title(stoneName)
	stoneKey := fmt.Sprintf("stone_%s", stoneTitle)
	res, err := rds.RedisGet(rds.RC, stoneKey)
	if err != nil {
		tools.WriteErr(err)
		messageSend(s, "Error retrieving Spirit Stone data!")
	} else {
		messageSend(s, res)
	}
}

func myTeamRouteSS(s *dg.Session, sender string, url string) {
	if url == "GET" {
		// redis get persons team image
		link, err := rds.RedisGet(rds.RC, sender)
		if err != nil {
			tools.WriteErr(err)
			messageSend(s, "You have not set your team!")
		} else {
			messageSend(s, link)
		}
	} else {
		// else set the url
		capsName := strings.Title(sender)
		ok := rds.RedisSet(rds.RC, capsName, url)
		if !ok {
			tools.WriteLog("Error: myTeamRouteSS() redisSet() Fail")
			messageSend(s, "Something went wrong, alert the Master Llama!")
		} else {
			messageSend(s, "Team set!")
		}

	}
}

func getTeamRouteSS(s *dg.Session, user string) {
	// redis get persons team image
	capName := strings.Title(user)
	link, err := rds.RedisGet(rds.RC, capName)
	if err != nil {
		tools.WriteErr(err)
		messageSend(s, "That person has no team set!")
	} else {
		messageSend(s, link)
	}
}

/*
	Soccer Spirits Tools
*/
// Search for the highest evolution of a player, starting at 3(EE), 2(E), 1(Base)
func highestEvoSS(playerName string) string {
	// Start at "_3"
	// HGETALL to see if entry exists, if not decrement, repeat
	finalForm := playerName + "_3"
	for i := 3; i > 0; i-- {
		lookupKey := strings.Title(fmt.Sprintf("%s_%d", playerName, i))
		exists, err := rds.RC.Exists(lookupKey).Result()
		if err != nil {
			tools.WriteErr(err)
		} else if exists {
			finalForm = lookupKey
			break
		} else {
			continue
		}
	}

	return finalForm
}
