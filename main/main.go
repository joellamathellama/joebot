package main

import (
	"flag"
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/api"
	"joebot/rds"
	"joebot/routes"
	"joebot/tools"
)

var (
	Token string
	BotID string
	err   error
)

/* Initialization for Discord bot, DB, redis, etc. */
func init() {
	// Init DB
	// initDB()

	tools.WriteLog("Bot Initializing")
	// Set flag variables
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()

	// Initiate redis
	fmt.Println("Init Redis. Expect no panic")
	rds.RedisInit()
	// Flush Redis
	// fmt.Println("Flushing ALL Keys in ALL Databases")
	// rds.RC.FlushAll()
	// Test redis Set & Get
	fmt.Println("Redis Set & Get test. Expect no panic")
	ok := rds.RedisSet(rds.RC, "redis_test_key", "redis_test_value")
	if !ok {
		tools.WriteLog("Redis Set Failed")
		fmt.Println("Redis Set Failed")
	}
	_, err = rds.RedisGet(rds.RC, "redis_test_key")
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Redis Get Failed")
	}
	// Test invalid query
	fmt.Println("Redis invalid query test. Expect no panic")
	rds.RedisGet(rds.RC, "nope")

	// Ssherder API call(s)
	// Called once on init and stored into Redis
	ssherderApis()

	// Create map of quick responses
	routes.BotResInit()
}

/* calls functions that GET Soccer Spirits player data */
func ssherderApis() {
	api.GetSkills()
	api.GetPlayers()
	api.GetStones()
}

/* whenReady sets current user's(the bot) status */
func whenReady(s *dg.Session, event *dg.Ready) {
	// Set the playing status.
	if err = s.UpdateStatus(0, "~joebot help"); err != nil {
		tools.WriteErr(err)
		fmt.Println(err)
	}
}

func main() {
	// Create a new Discord session using the bot token
	dg, err := dg.New(Token)
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Update status on ready
	dg.AddHandler(whenReady)

	// Register messageRoutes as a callback for the messageRoutes events.
	dg.AddHandler(routes.MessageRoutes)

	// Open the websocket and begin listening.
	if err = dg.Open(); err != nil {
		tools.WriteErr(err)
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
