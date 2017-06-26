package main

import (
	"flag"
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/alarms"
	"joebot/api"
	"joebot/rds"
	"joebot/routes"
	t "joebot/tools"
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

	t.WriteLog("Bot Initializing")
	// Set flag variables
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()

	// Initiate redis
	fmt.Println("Init Redis. Expect no panic")
	rds.RedisInit()
	// Flush Redis
	// fmt.Println("Flushing Redis")
	// rds.RC.FlushAll()
	// Test redis Set & Get
	fmt.Println("Redis Set & Get test. Expect no panic")
	ok := rds.RedisSet(rds.RC, "redis_test_key", "redis_test_value")
	if !ok {
		t.WriteLog("Redis Set Failed")
		fmt.Println("Redis Set Failed")
	}
	_, err = rds.RedisGet(rds.RC, "redis_test_key")
	if err != nil {
		t.WriteErr(err)
		fmt.Println("Redis Get Failed")
	}
	// Test invalid query
	fmt.Println("Redis invalid query test. Expect no panic")
	rds.RedisGet(rds.RC, "nope")

	// Ssherder API call(s)
	// Called once on init and stored into Redis
	ssherderApis()

	// JSON to Redis
	dwuToRedis()
	api.LocalizationToRedis()
	api.PilotsToRedis()
	api.Test1()
	api.Test2()

	// Create map of quick responses
	routes.BotResInit()

	// Create lists
	routes.CreateAlarmList()
}

/* api calls */
func ssherderApis() {
	t.WriteLog("ssherderApis()")
	api.GetSkills()
	api.GetPlayers()
	api.GetStones()
}

/* dwu json to redis */
func dwuToRedis() {
	t.WriteLog("dwuToRedis()")

	api.OfficersToRedis()

	list := [5]string{"wei", "shu", "wu", "other", "jin"}
	for ii := 0; ii < len(list); ii++ {
		api.PassivesToRedis(list[ii])
	}
}

/* whenReady sets current user's(the bot) status */
func whenReady(s *dg.Session, event *dg.Ready) {
	// Set the playing status.
	t.WriteLog("whenReady()")

	bStatus := "~help"
	if err = s.UpdateStatus(0, bStatus); err != nil {
		t.WriteErr(err)
		fmt.Println(err)
	} else {
		bS := fmt.Sprintf("Update Status: %s", bStatus)
		t.WriteLog(bS)
	}

	// Start Tickers
	alarms.AlarmGKShootout(s)
}

func main() {
	// Create a new Discord session using the bot token
	t.WriteLog("main()")
	dg, err := dg.New(Token)
	if err != nil {
		t.WriteErr(err)
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		t.WriteErr(err)
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
		t.WriteErr(err)
		fmt.Println("error opening connection,", err)
		return
	} else {
		t.WriteLog("Discord websocket connection open")
	}

	wLog := "Bot is now running.  Press CTRL-C to exit."
	t.WriteLog(wLog)
	fmt.Println(wLog)

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
