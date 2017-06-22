package routes

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/rds"
	"joebot/tools"
	"strings"
)

func skillsRouteGK(s *dg.Session, pilotName string) (res string) {
	// rds.RedisGet(key)
	// if else error
	pilotLow := strings.ToLower(pilotName)
	pk := fmt.Sprintf("gk_ps_%s", pilotLow)
	res, err := rds.RedisGet(rds.RC, pk)
	if err != nil {
		tools.WriteErr(err)
		res = "Error retrieving Pilot Skills!"
	}
	return
}
