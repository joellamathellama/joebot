package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	autoRes map[string]string
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
	autoRes = make(map[string]string)
	// Fill it up
	autoResInit(autoRes)
}

func autoResInit(m map[string]string) {
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
	} else if len(c) < 8 || c[:8] != "~joebot " {
		return
	}

	if stringInSlice(c[8:], cmdList) {// cmdList defined in func autoResInit
		messageSend(s, cID, autoRes[c[8:]])
	} else if c[8:14] == "search" {
		res, _ := redisClient.HGet(c[15:], "Story").Result()
		// fmt.Println(res)
		messageSend(s, cID, res)
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
