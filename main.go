package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Globals
/*
	cID = current channel ID
	cmdResList = map of commands and the corresponding responses
*/
var (
	cID string
	cmdResList map[string]string
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
	// Called once on init and stored into Redis
	getChars()

	// Create map of quick responses
	botResInit()
}

func botResInit() {
	cmdResList = make(map[string]string)

	// Fill it up
	cmdResList["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	cmdResList["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	cmdResList["reddit"] = "http://reddit.com/r/soccerspirits"
}

func messageSend(s *discordgo.Session, m string) {
	if _, err = s.ChannelMessageSend(cID, m); err != nil {
		// fmt.Println("Error - s.ChannelMessageSend: ", err)
		panic(err)
	}
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageRoutes(s *discordgo.Session, m *discordgo.MessageCreate) {
	// c = users message, cID = channel ID
	c := m.Content
	cID = m.ChannelID

	// Ignore all messages created by the bot itself and anything short of "~joebot "
	if m.Author.ID == BotID {
		return
	} else if len(c) < 8 || regexpMatch("^(?i)(~Joebot)[ ]", c[:8]) != true {
		return
	}

	// cmdList = valid quick commands
	cmdList := []string{"ourteams", "apoc", "reddit"}

	// ROUTES
	if stringInSlice(c[8:], cmdList) {
		messageSend(s, cmdResList[c[8:]])
	} else if regexpMatch("(?i)(Story)[ ][a-zA-Z0-9]", c[8:]) {
		storyRoute(s, c[14:])
	} else if regexpMatch("(?i)(Stones)[ ][a-zA-Z0-9]", c[8:]) {
		stonesRoute(s, c[15:])
	}else {
		messageSend(s, "Enter a valid command")
	}
}

func storyRoute(s *discordgo.Session, playerName string) {
	res, err := redisClient.HGet(strings.Title(playerName), "Story").Result()
	if err != nil {
		messageSend(s, "Enter a valid command")
	} else {
		messageSend(s, res)
	}
}

func stonesRoute(s *discordgo.Session, playerName string) {
	res, err := redisClient.HGet(strings.Title(playerName), "Stones").Result()
	if err != nil {
		messageSend(s, "Enter a valid command")
	} else {
		messageSend(s, res)
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

	// Register messageRoutes as a callback for the messageRoutes events.
	dg.AddHandler(messageRoutes)

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
