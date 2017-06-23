package routes

import (
	"fmt"
	// "reflect"
	dg "github.com/bwmarrin/discordgo"
	"joebot/tools"
	"strings"
)

/*
	cID = current channel ID
	cmdResList = map of commands and the corresponding responses
*/
var (
	cmdResList map[string]string
	BotID      string
	err        error
)

func SendMessage(s *dg.Session, cID string, m string) {
	// fmt.Printf("@@@: %s, %s", cID, m)
	if _, err = s.ChannelMessageSend(cID, m); err != nil {
		tools.WriteErr(err)
		fmt.Println(err)
	}
}

// Quick bot responses
func BotResInit() {
	cmdResList = make(map[string]string)

	// Fill it up
	cmdResList["ourteams"] = "https://docs.google.com/spreadsheets/d/1ykMKW64o71OSfOEtx-iIa25jSZCFVRcZQ73ErXEoFpc/edit#gid=0"
	cmdResList["apoc"] = "http://soccerspirits.freeforums.net/thread/69/guide-apocalypse-player-tier-list"
	cmdResList["redditss"] = "http://reddit.com/r/soccerspirits"
	cmdResList["redditdwu"] = "http://reddit.com/r/dwunleashed"
	cmdResList["teamwork"] = "https://docs.google.com/spreadsheets/d/1x0Q4vUk_V3wUwzM5XR_66xytSbapoSFm_cHR9PYIERs/htmlview?sle=true#"
	cmdResList["chains"] = "https://ssherder.com/characters/#"
	cmdResList["help"] = "Shoutout to ssherder.com, api.lootbox.eu/documentation#/ and gkgirls.info.gf/\n\n" +
		"*General Commands:*\n**Write my own Note:** '~mynote <Text>'\n" +
		"**Read others Note:** '~note <Discord Name>'\n" +
		"**Alarms:** '~alarm <Name>'\n\n" +
		"*Overwatch Commands:(Lootbox seems to be down for now)*\n" +
		"**Lookup PC Profile:** '~PCprofile <Bnet Tag>'\n" +
		"**Lookup PC Stats:** '~PCstats <Bnet Tag>'\n" +
		"**Lookup PS:** Same thing, except '~PSprofile, ~PSstats'\n" +
		"**Lookup Xbox:** Same thing, except '~Xprofile, ~Xstats'\n\n" +
		"*Soccer Spirits Commands:*\n**Lookup player info:** '~sstory, ~sstone, ~sslots, ~ssherder or ~sskills <Name>'\n" +
		"**Quick links:** '~ourteams', '~apoc', '~reddit'\n\n" +
		"*Dynasty Warriors Unleashed Commands:*\n" +
		"**Lookup Officer Legendary Passives:** '~dwup <Name>'\n" +
		"**Lookup Officer Stats:** '~dwus <Name>'\n\n" +
		"*Goddess Kiss Commands:*\n**Lookup Pilot Skills:** '~gskills <Name>'\n\n" +
		"*Everything is case *insensitive!*(Except Bnet Tags)"
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageRoutes(s *dg.Session, m *dg.MessageCreate) {
	// Contents
	c := m.Content // full message sent by user

	// Meta
	cID := m.ChannelID
	sender := m.Author.Username
	// Ignore all messages created by the bot itself and anything short of "~"
	if m.Author.ID == BotID {
		return
	} else if len(c) < 2 || c[0:1] != "~" {
		return
	}

	// split message by command and arguments
	cSplit := strings.Split(c[1:], " ") // ["command", ..., ...]
	cc := cSplit[0]                     // "command"
	cl := len(cc) + 2
	cmdArgs := ""

	if len(cmdResList[cc]) != 0 { // if quick command
		SendMessage(s, cID, cmdResList[cc])
		return
	} else if len(cSplit) >= 2 {
		cmdArgs = c[cl:]
	}

	/*
		ROUTES
	*/
	res := ""
	switch ccLow := strings.ToLower(cc); ccLow {
	/* General */
	case "mynote":
		if len(cmdArgs) > 0 {
			res = myTeamRouteSS(s, sender, cmdArgs)
		} else {
			res = myTeamRouteSS(s, sender, "GET")
		}
	case "note":
		res = getTeamRouteSS(s, cmdArgs)
	case "setalarm":
		res = setAlarm(s, cID, cmdArgs)
	case "removealarm":
		res = removeAlarm(s, cID, cmdArgs)
	/* Soccer Spirits */
	case "sstory":
		res = storyRouteSS(s, cmdArgs)
	case "sslots":
		res = slotesRouteSS(s, cmdArgs)
	case "ssherder":
		res = ssherderRouteSS(s, cmdArgs)
	case "sskills":
		res = skillsRouteSS(s, cmdArgs)
	case "sstone":
		res = stoneRouteSS(s, cmdArgs)
		/* Overwatch */
	case "pcprofile":
		res = profileRouteOW(s, cmdArgs, "pc")
	case "pcstats":
		res = statsRouteOW(s, cmdArgs, "pc")
	case "psprofile":
		res = profileRouteOW(s, cmdArgs, "psn")
	case "psstats":
		res = statsRouteOW(s, cmdArgs, "psn")
	case "xprofile":
		res = profileRouteOW(s, cmdArgs, "xbl")
	case "xstats":
		res = statsRouteOW(s, cmdArgs, "xbl")
		/* Dynasty Warriors Unleashed */
	case "dwup":
		res = passiveRouteDWU(s, cmdArgs)
	case "dwus":
		res = officerRouteDWU(s, cmdArgs)
		/* Goddess Kiss */
	case "gskills":
		res = skillsRouteGK(s, cmdArgs)
	default:
		res = "Enter a valid command"
	}

	SendMessage(s, cID, res)
	return
}
