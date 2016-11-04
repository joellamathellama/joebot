package main

import (
	"fmt"
	// "reflect"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func messageSend(s *dg.Session, m string) {
	if _, err = s.ChannelMessageSend(cID, m); err != nil {
		writeErr(err)
		fmt.Println(err)
		checkError(err)
	}
}

// Quick bot responses
func botResInit() {
	cmdResList = make(map[string]string)

	// Fill it up
	cmdResList["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	cmdResList["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	cmdResList["reddit"] = "http://reddit.com/r/soccerspirits"
	cmdResList["help"] = "*Overwatch Commands:*\n**Lookup PC Profile:** '~joebot PCprofile <Battlenet Tag>' (Ex. ~joebot pcprofile joellama#1114)\n**Lookup PC Stats:** '~joebot PCstats <Battlenet Tag>' (Ex. ~joebot pcstats joellama#1114)\n**Lookup PS:** Same thing, except 'PSprofile/PSstats'\n**Lookup Xbox:** Same thing, except 'Xprofile/Xstats'\n\n*Soccer Spirits Commands:*\n**Lookup player info:** '~joebot Story, Stones, Ssherder or Skills <Player Name>' (Ex. ~joebot stats Griffith)\n**Quick links:** 'ourteams', 'apoc', 'reddit' (Ex. ~joebot apoc)\n\n*Everything is case *insensitive!*(Except Bnet Tags)"
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageRoutes(s *dg.Session, m *dg.MessageCreate) {
	// c = users message, cID = channel ID
	c := m.Content
	cID = m.ChannelID

	// Ignore all messages created by the bot itself and anything short of "~joebot "
	if m.Author.ID == BotID {
		return
	} else if len(c) < 8 || regexpMatch("^(?i)(~Joebot)[ ]", c[:8]) != true {
		return
	}

	/*
		ROUTES
		if regexpMatch(REGEX, command)
			routeFunction()
		else
			Send message "command not found"
	*/
	if len(cmdResList[c[8:]]) != 0 {
		messageSend(s, cmdResList[c[8:]])
		/* SOCCER SPIRITS */
	} else if regexpMatch("(?i)(Story)[ ][a-zA-Z0-9]", c[8:]) {
		storyRouteSS(s, c[14:])
	} else if regexpMatch("(?i)(Stones)[ ][a-zA-Z0-9]", c[8:]) {
		stonesRouteSS(s, c[15:])
	} else if regexpMatch("(?i)(Ssherder)[ ][a-zA-Z0-9]", c[8:]) {
		ssherderRouteSS(s, c[17:])
	} else if regexpMatch("(?i)(Skills)[ ][a-zA-Z0-9]", c[8:]) {
		skillsRouteSS(s, c[15:])
		/* OVERWATCH */
	} else if regexpMatch("(?i)(PCprofile)[ ][a-zA-Z0-9]", c[8:]) {
		profileRouteOW(s, c[18:], "pc")
	} else if regexpMatch("(?i)(PCstats)[ ][a-zA-Z0-9]", c[8:]) {
		statsRouteOW(s, c[16:], "pc")
	} else if regexpMatch("(?i)(PSprofile)[ ][a-zA-Z0-9]", c[8:]) {
		profileRouteOW(s, c[18:], "psn")
	} else if regexpMatch("(?i)(PSstats)[ ][a-zA-Z0-9]", c[8:]) {
		statsRouteOW(s, c[16:], "psn")
	} else if regexpMatch("(?i)(Xprofile)[ ][a-zA-Z0-9]", c[8:]) {
		profileRouteOW(s, c[17:], "xbl")
	} else if regexpMatch("(?i)(Xstats)[ ][a-zA-Z0-9]", c[8:]) {
		statsRouteOW(s, c[15:], "xbl")
	} else {
		messageSend(s, "Enter a valid command")
	}
}

/*
	SOCCER SPIRITS ROUTES
*/

// Search for the highest evolution of a player, starting at 3(EE), 2(E), 1(Base)
func highestEvoSS(playerName string) string {
	// Start at "_3"
	// HGETALL to see if entry exists, if not decrement, repeat
	finalForm := playerName + "_3"
	for i := 3; i > 0; i-- {
		lookupKey := strings.Title(fmt.Sprintf("%s_%d", playerName, i))
		exists, err := rc.Exists(lookupKey).Result()
		if err != nil {
			writeErr(err)
		} else if exists {
			finalForm = lookupKey
			break
		} else {
			continue
		}
	}

	return finalForm
}

func storyRouteSS(s *dg.Session, playerName string) {
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rc.HGet(lookupKey, "Story").Result()
	if err != nil {
		writeErr(err)
		messageSend(s, "Player's story not found! Try, idk, typing it correctly?")
	} else {
		messageSend(s, res)
	}
}

func stonesRouteSS(s *dg.Session, playerName string) {
	lookupKey := highestEvoSS(playerName)
	// fmt.Println(lookupKey)
	res, err := rc.HGet(lookupKey, "Stones").Result()
	if err != nil {
		writeErr(err)
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
	res, err := rc.HGet(lookupKey, "ID").Result()
	if err != nil {
		writeErr(err)
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
	res, err := rc.HGet(lookupKey, "Skills").Result()
	if err != nil {
		writeErr(err)
		messageSend(s, "Player's skills not found. Sharpen your typing skills first...")
	} else {
		messageSend(s, res)
	}
}

/* OVERWATCH ROUTES */

func profileRouteOW(s *dg.Session, playerName string, platform string) {
	// replace # with - and call getPlayerStats
	fmtName := strings.Replace(playerName, "#", "-", -1)

	// Look it up in redis, if exit, return info, if not, continue
	playerHash := fmt.Sprintf("%s%s", fmtName, platform)
	res, err := rc.HGet(playerHash, "profile").Result()
	if err != nil {
		writeErr(err)
	} else {
		messageSend(s, res)
		return
	}

	messageSend(s, "This may take a few seconds...")

	playerProfile := getPlayerProfile(fmtName, platform)

	messageSend(s, playerProfile)
}

func statsRouteOW(s *dg.Session, playerName string, platform string) {
	fmtName := strings.Replace(playerName, "#", "-", -1)

	// Look it up in redis, if exit, return info, if not, continue
	playerHash := fmt.Sprintf("%s%s", fmtName, platform)
	res, err := rc.HGet(playerHash, "stats").Result()
	if err != nil {
		writeErr(err)
	} else {
		messageSend(s, res)
		return
	}

	messageSend(s, "This may take few seconds...")

	playerStats := getPlayerStats(fmtName, platform)

	messageSend(s, playerStats)
}
