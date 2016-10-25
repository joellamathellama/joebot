package main

import (
	"fmt"
	"flag"

	dg "github.com/bwmarrin/discordgo"
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
	redisSet(rc, "redis_test_key", "redis_test_value")
	redisGet(rc, "redis_test_key")
	// Test invalid query
	redisGet(rc, "nope")

	// Ssherder API call(s)
	// Called once on init and stored into Redis
	ssherderApis()

	// Create map of quick responses
	botResInit()
}

func ssherderApis() {
	getSkills()
	defer getPlayers()
}

func botResInit() {
	cmdResList = make(map[string]string)

	// Fill it up
	cmdResList["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	cmdResList["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	cmdResList["reddit"] = "http://reddit.com/r/soccerspirits"
}

func messageSend(s *dg.Session, m string) {
	if _, err = s.ChannelMessageSend(cID, m); err != nil {
		// fmt.Println("Error - s.ChannelMessageSend: ", err)
		panic(err)
	}
}

func main() {
	// Create a new Discord session using the bot token
	dg, err := dg.New(Token)
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
