package main

import (
	// "fmt"
	"strings"
	// "reflect"

	dg "github.com/bwmarrin/discordgo"
)

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
	} else if regexpMatch("(?i)(Skills)[ ][a-zA-Z0-9]", c[8:]) {
		skillsRoute(s, c[15:])
	}else {
		messageSend(s, "Enter a valid command")
	}
}

func storyRoute(s *dg.Session, playerName string) {
	res, err := rc.HGet(strings.Title(playerName), "Story").Result()
	if err != nil {
		messageSend(s, "Enter a valid command")
	} else {
		messageSend(s, res)
	}
}

func stonesRoute(s *dg.Session, playerName string) {
	res, err := rc.HGet(strings.Title(playerName), "Stones").Result()
	if err != nil {
		messageSend(s, "Enter a valid command")
	} else {
		messageSend(s, res)
	}
}

func skillsRoute(s *dg.Session, playerName string) {
	playerKey := strings.Title(playerName) + "_skills"

	// fmt.Println(rc.LLen(playerKey))

	res, err := rc.LPop(playerKey).Result()
	if err != nil {
		messageSend(s, "Enter a valid command")
	} else {
		messageSend(s, res)
	}
}