package main

import (
	"fmt"
	"strings"
	// "reflect"

	dg "github.com/bwmarrin/discordgo"
)

func messageSend(s *dg.Session, m string) {
	if _, err = s.ChannelMessageSend(cID, m); err != nil {
		// fmt.Println("Error - s.ChannelMessageSend: ", err)
		panic(err)
	}
}

// Quick bot responses
func botResInit() {
	cmdResList = make(map[string]string)

	// Fill it up
	cmdResList["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	cmdResList["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	cmdResList["reddit"] = "http://reddit.com/r/soccerspirits"
	cmdResList["help"] = "*Overwatch Commands:*\n**Lookup PC stats:** '~joebot pcstats <Battlenet Tag>' (Ex. ~joebot pcstats joellama#1114)\n\n*Soccer Spirits Commands:*\n**Lookup player info:** '~joebot Story, Stones, Ssherder or Skills <Player Name>' (Ex. ~joebot stats Griffith)\n**Quick links:** 'ourteams', 'apoc', 'reddit' (Ex. ~joebot apoc)"
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

	// ROUTES
	if len(cmdResList[c[8:]]) != 0 {
		messageSend(s, cmdResList[c[8:]])
	} else if regexpMatch("(?i)(Story)[ ][a-zA-Z0-9]", c[8:]) {
		storyRoute(s, c[14:])
	} else if regexpMatch("(?i)(Stones)[ ][a-zA-Z0-9]", c[8:]) {
		stonesRoute(s, c[15:])
	} else if regexpMatch("(?i)(Ssherder)[ ][a-zA-Z0-9]", c[8:]) {
		ssherderRoute(s, c[17:])
	} else if regexpMatch("(?i)(Skills)[ ][a-zA-Z0-9]", c[8:]) {
		skillsRoute(s, c[15:])
	} else if regexpMatch("(?i)(pcstats)[ ][a-zA-Z0-9]", c[8:]) {
		statsRoute(s, c[16:])
	} else {
		messageSend(s, "Enter a valid command")
	}
}

/*
	the appended "_3" referrs to a players third evolution, which is all anyone cares about
*/
func storyRoute(s *dg.Session, playerName string) {
	lookupKey := strings.Title(playerName) + "_3"
	fmt.Println(lookupKey)
	res, err := rc.HGet(lookupKey, "Story").Result()
	if err != nil {
		messageSend(s, "Player's story not found! Try, idk, typing it correctly?")
	} else {
		messageSend(s, res)
	}
}

func stonesRoute(s *dg.Session, playerName string) {
	lookupKey := strings.Title(playerName) + "_3"
	fmt.Println(lookupKey)
	res, err := rc.HGet(lookupKey, "Stones").Result()
	if err != nil {
		messageSend(s, "Player's stones not found. Prolly cause you're stoned...")
	} else {
		messageSend(s, res)
	}
}

func ssherderRoute(s *dg.Session, playerName string) {
	// https://ssherder.com/characters/ID/
	// lookup player ID, add to URL, send message
	lookupKey := strings.Title(playerName) + "_3"
	fmt.Println(lookupKey)
	res, err := rc.HGet(lookupKey, "ID").Result()
	if err != nil {
		messageSend(s, "Who?!")
	} else {
		messageSend(s, "https://ssherder.com/characters/"+res)
	}
}

func skillsRoute(s *dg.Session, playerName string) {
	// https://ssherder.com/characters/ID/
	// lookup player ID, add to URL, send message
	lookupKey := strings.Title(playerName) + "_3"
	fmt.Println(lookupKey)
	res, err := rc.HGet(lookupKey, "Skills").Result()
	if err != nil {
		messageSend(s, "Player's skills not found. Sharpen your typing skills first...")
	} else {
		messageSend(s, res)
	}
}

func statsRoute(s *dg.Session, playerName string) {
	// replace # with - and call getPlayerStats
	fmtName := strings.Replace(playerName, "#", "-", -1)
	playerStats := getPlayerStats(fmtName)
	// playerStats := getPlayerStats("jawnkeem-1982")
	messageSend(s, playerStats)
}
