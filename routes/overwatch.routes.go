package routes

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/api"
	"joebot/rds"
	"joebot/tools"
	"strings"
)

func profileRouteOW(s *dg.Session, playerName string, platform string) {
	// replace # with - and call getPlayerStats
	fmtName := strings.Replace(playerName, "#", "-", -1)

	// Look it up in redis, if exit, return info, if not, continue
	playerHash := fmt.Sprintf("%s%s", fmtName, platform)
	res, err := rds.RC.HGet(playerHash, "profile").Result()
	if err != nil {
		tools.WriteErr(err)
	} else {
		messageSend(s, res)
		return
	}

	messageSend(s, "This may take a few seconds...")

	playerProfile := api.GetPlayerProfile(fmtName, platform)

	messageSend(s, playerProfile)
}

func statsRouteOW(s *dg.Session, playerName string, platform string) {
	fmtName := strings.Replace(playerName, "#", "-", -1)

	// Look it up in redis, if exit, return info, if not, continue
	playerHash := fmt.Sprintf("%s%s", fmtName, platform)
	res, err := rds.RC.HGet(playerHash, "stats").Result()
	if err != nil {
		tools.WriteErr(err)
	} else {
		messageSend(s, res)
		return
	}

	messageSend(s, "This may take few seconds...")

	playerStats := api.GetPlayerStats(fmtName, platform)

	messageSend(s, playerStats)
}
