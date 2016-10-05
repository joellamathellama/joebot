package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"gopkg.in/redis.v4"
)

// Variables used for command line parameters
var (
	redisClient *redis.Client
	autoRes map[string]string
	Token string
	BotID string
	err error
)

func init() {
	// set flag variables
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()

	// Initiate redis(Not used, but ready for use)
	redisInit()

	// Switch this from a map to redis cache for fun? Maybe...
	// create empty map for auto responses
	autoRes = make(map[string]string)
	// fill it up
	fillAutoRes(autoRes)
}

func fillAutoRes(m map[string]string) {
	m["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	m["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	m["reddit"] = "http://reddit.com/r/soccerspirits"
}

// Connect to default port
func redisInit() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	// redist test: WORKING
	// redisSet(redisClient, "test key", "test string")
	// redisGet(redisClient, "test key")
}

func redisSet(c *redis.Client, key string, value string) {
	err := c.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func redisGet(c *redis.Client, key string) {
	val, err := c.Get(key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(key, val)
}

func messageSend(s *discordgo.Session, c, m string) {
	_, err = s.ChannelMessageSend(c, m)
	if err != nil {
		// fmt.Println("Error - s.ChannelMessageSend: ", err)
		panic(err)
	}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// c = users message, cID = channel ID
	c, cID := m.Content, m.ChannelID

	// Ignore all messages created by the bot itself and anything short of "~joebot"
	if m.Author.ID == BotID {
		return
	} else if len(c) < 9 || c[:8] != "~joebot " {
		return
	}

	// cmd = anything after "~joebot "
	cmd := c[8:]
	cmdList := []string{"ourteams", "apoc", "reddit"}
	if stringInSlice(cmd, cmdList) {
		messageSend(s, cID, autoRes[cmd])
	} else {
		messageSend(s, cID, "Enter a valid command")
	}
}

func main() {
	// Create a new Discord session using the provided login information.
	// dg, err := discordgo.New(Email, Password, Token)
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
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
