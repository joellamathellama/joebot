package main

import (
	"flag"
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

// Globals
/*
	cID = current channel ID
	cmdResList = map of commands and the corresponding responses
*/
var (
	cID        string
	cmdResList map[string]string
	Token      string
	BotID      string
	err        error
)

func init() {
	// Set flag variables
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()

	// Initiate redis
	fmt.Println("Init Redis. Expect: No response")
	redisInit()
	// Flush Redis
	fmt.Println("Flushing ALL Keys in ALL Databases")
	rc.FlushAll()
	// Test redis Set & Get
	fmt.Println("Redis Set & Get test. Expect: No response")
	redisSet(rc, "redis_test_key", "redis_test_value")
	redisGet(rc, "redis_test_key")
	// Test invalid query
	fmt.Println("Redis invalid query test. Expect: 'Invalid Key'")
	redisGet(rc, "nope")

	// Ssherder API call(s)
	// Called once on init and stored into Redis
	ssherderApis()

	// Create map of quick responses
	botResInit()
}

func ssherderApis() {
	getSkills()
	getPlayers()
}

func whenReady(s *dg.Session, event *dg.Ready) {
	// Set the playing status.
	if err = s.UpdateStatus(0, "Type: '~joebot help'"); err != nil {
		fmt.Println(err)
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

	// Update status on ready
	dg.AddHandler(whenReady)

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
