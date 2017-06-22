package routes

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/rds"
	t "joebot/tools"
	"strings"
)

func passiveRouteDWU(s *dg.Session, playerName string) (res string) {
	playerName = strings.ToLower(playerName)
	lookupKey := fmt.Sprintf("passive_%s", playerName)

	res, err := rds.RedisGet(rds.RC, lookupKey)
	if err != nil {
		t.WriteErr(err)
		res = "Player's passives not found! Make sure it's spelled correctly!"
	}
	return
}

func officerRouteDWU(s *dg.Session, playerName string) (res string) {
	playerName = strings.ToLower(playerName)
	lookupKey := fmt.Sprintf("officer_%s", playerName)

	res, err := rds.RedisGet(rds.RC, lookupKey)
	if err != nil {
		t.WriteErr(err)
		res = "Officer not found! Make sure it's spelled correctly!"
	}
	return
}
