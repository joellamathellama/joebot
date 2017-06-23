package routes

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/api"
	"joebot/rds"
	"joebot/tools"
	"strings"
)

func profileRouteOW(s *dg.Session, playerName string, platform string) (res string) {
	// replace # with - and call getPlayerStats
	fmtName := strings.Replace(playerName, "#", "-", -1)

	// Look it up in redis, if exit, return info, if not, continue
	playerHash := fmt.Sprintf("%s%s", fmtName, platform)
	res, err := rds.RC.HGet(playerHash, "profile").Result()
	if err != nil {
		tools.WriteErr(err)
	} else {
		return
	}

	playerProfile := api.GetPlayerProfile(fmtName, platform)
	res = playerProfile
	return
}

func statsRouteOW(s *dg.Session, playerName string, platform string) (res string) {
	fmtName := strings.Replace(playerName, "#", "-", -1)

	// Look it up in redis, if exit, return info, if not, continue
	playerHash := fmt.Sprintf("%s%s", fmtName, platform)
	res, err := rds.RC.HGet(playerHash, "stats").Result()
	if err != nil {
		tools.WriteErr(err)
	} else {
		return
	}

	playerStats := api.GetPlayerStats(fmtName, platform)
	res = playerStats
	return
}
