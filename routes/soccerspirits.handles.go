package routes

import (
	"fmt"
	"joebot/rds"
	"joebot/tools"
	"strings"
)

func storyRouteSS(playerName string) (res string) {
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "Story").Result()
	if err != nil {
		tools.WriteErr(err)
		res = "Player's story not found! Try, idk, typing it correctly?"
	}
	return
}

func slotesRouteSS(playerName string) (res string) {
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "Stones").Result()
	if err != nil {
		tools.WriteErr(err)
		res = "Player's stones not found. Prolly cause you're stoned..."
	}
	return
}

func ssherderRouteSS(playerName string) (res string) {
	// https://ssherder.com/characters/ID/
	// lookup player ID, add to URL, send message
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "ID").Result()
	if err != nil {
		tools.WriteErr(err)
		res = "Who?!"
	} else {
		res = fmt.Sprintf("https://ssherder.com/characters/%s", res)
	}
	return
}

func skillsRouteSS(playerName string) (res string) {
	// https://ssherder.com/characters/ID/
	// lookup player ID, add to URL, send message
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rds.RC.HGet(lookupKey, "Skills").Result()
	if err != nil {
		tools.WriteErr(err)
		res = "Player's skills not found. Sharpen your typing skills first..."
	}
	return
}

func stoneRouteSS(stoneName string) (res string) {
	// rds.RedisGet(key)
	// if else error
	stoneTitle := strings.ToLower(stoneName)
	stoneKey := fmt.Sprintf("stone_%s", stoneTitle)
	res, err := rds.RedisGet(rds.RC, stoneKey)
	if err != nil {
		tools.WriteErr(err)
		res = "Error retrieving Spirit Stone data!"
	}
	return
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
		lookupKey := strings.ToLower(fmt.Sprintf("%s_%d", playerName, i))
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
