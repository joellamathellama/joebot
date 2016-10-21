package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	cmdResList map[string]string
	cmdList []string
	Token string
	BotID string
	err error
)

func init() {
	// Set flag variables
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()

	// Initiate redis
	redisInit()
	// Test redis Set & Get
	redisSet(redisClient, "redis_test_key", "redis_test_value")
	redisGet(redisClient, "redis_test_key")
	// Test invalid query
	redisGet(redisClient, "nope")

	// Ssherder API call(s)
	getChars()

	// Create empty map for auto responses
	cmdResList = make(map[string]string)
	// Fill it up
	botResInit(cmdResList)
}

func botResInit(m map[string]string) {
	// Slice of valid commands used in func messageCreate
	cmdList = []string{"ourteams", "apoc", "reddit"}

	// Command responses
	m["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	m["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	m["reddit"] = "http://reddit.com/r/soccerspirits"
}

func messageSend(s *discordgo.Session, c, m string) {
	// fmt.Println("Channel id: ", c)
	if _, err = s.ChannelMessageSend(c, m); err != nil {
		// fmt.Println("Error - s.ChannelMessageSend: ", err)
		panic(err)
	}
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// c = users message, cID = channel ID
	c, cID := m.Content, m.ChannelID

	// Ignore all messages created by the bot itself and anything short of "~joebot "
	if m.Author.ID == BotID {
		return
	} else if len(c) < 8 || regexpMatch("^(?i)(~Joebot)[ ]", c[:8]) != true {
		return
	}

	// condition to match "~joebot", " ", {command}, {name}
	// take a diff path depending if there is a third argument or not

	// I'll think of something better...
	if stringInSlice(c[8:], cmdList) {// cmdList defined in func botResInit
		messageSend(s, cID, cmdResList[c[8:]])
	} else if regexpMatch("(?i)(Story)[ ][a-zA-Z0-9]", c[8:]) {
		res, err := redisClient.HGet(strings.Title(c[14:]), "Story").Result()
		if err != nil {
			messageSend(s, cID, "Enter a valid command")
		} else {
			messageSend(s, cID, res)
		}
	} else if regexpMatch("(?i)(Stones)[ ][a-zA-Z0-9]", c[8:]) {
		res, err := redisClient.HGet(strings.Title(c[15:]), "Stones").Result()
		if err != nil {
			messageSend(s, cID, "Enter a valid command")
		} else {
			messageSend(s, cID, res)
		}
	}else {
		messageSend(s, cID, "Enter a valid command")
	}
}

func main() {
	// Create a new Discord session using the bot token
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	if err = dg.Open(); err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
