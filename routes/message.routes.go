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
	cID        string
	cmdResList map[string]string
	BotID      string
	err        error
)

func messageSend(s *dg.Session, m string) {
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
	cmdResList["reddit"] = "http://reddit.com/r/soccerspirits"
	cmdResList["teamwork"] = "https://docs.google.com/spreadsheets/d/1x0Q4vUk_V3wUwzM5XR_66xytSbapoSFm_cHR9PYIERs/htmlview?sle=true#"
	cmdResList["chains"] = "https://ssherder.com/characters/#"
	cmdResList["help"] = "Shoutout to ssherder.com and api.lootbox.eu/documentation#/ for their APIs.\n\nTo talk to the bot, preface your message with '~joebot'\n\n*Overwatch Commands:*\n**Lookup PC Profile:** 'PCprofile <Battlenet Tag>' (Ex. ~joebot pcprofile joellama#1114)\n**Lookup PC Stats:** 'PCstats <Battlenet Tag>' (Ex. ~joebot pcstats joellama#1114)\n**Lookup PS:** Same thing, except 'PSprofile, PSstats'\n**Lookup Xbox:** Same thing, except 'Xprofile, Xstats'\n\n*Soccer Spirits Commands:*\n**Lookup player info:** 'Story, Stones, Ssherder or Skills <Player Name>' (Ex. ~joebot stats Griffith)\n**Quick links:** 'ourteams', 'apoc', 'reddit' (Ex. ~joebot apoc)\n**My Team:** 'myteam <TEAM or LINK>'(Ex. ~joebot myteam sharr/elaine/renee...)\n**Others Teams:** 'team <USERNAME>'(Ex. ~joebot team mazdarx13)\n\n*Everything is case *insensitive!*(Except Bnet Tags)"
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageRoutes(s *dg.Session, m *dg.MessageCreate) {
	// Contents
	c := m.Content // full message sent by user
	// Ignore all messages created by the bot itself and anything short of "~joebot "
	if m.Author.ID == BotID {
		return
	} else if len(c) < 8 || tools.RegexpMatch("^(?i)(~Joebot)", c[0:7]) != true {
		// fmt.Println("Not talking to bot")
		return
	}

	// split message by command and arguments
	cSplit := strings.Split(c, " ") // ["~joebot", "command", [...]]
	cc := cSplit[1]
	cl := len(cc) + 7
	if len(cSplit) > 2 { // extract argument
		cl = cl + 2
	} else {
		cl = cl + 1
	}
	cmdArgs := c[cl:]

	// Meta
	cID = m.ChannelID
	sender := m.Author.Username

	if len(cmdResList[cc]) != 0 { // if quick command
		messageSend(s, cmdResList[cc])
	}
	/*
		ROUTES
	*/
	switch ccLow := strings.ToLower(cc); ccLow {
	case "story":
		storyRouteSS(s, cmdArgs)
	case "slots":
		slotesRouteSS(s, cmdArgs)
	case "ssherder":
		ssherderRouteSS(s, cmdArgs)
	case "skills":
		skillsRouteSS(s, cmdArgs)
	case "stone":
		stoneRouteSS(s, cmdArgs)
	case "myteam":
		if len(cmdArgs) > 0 {
			myTeamRouteSS(s, sender, cmdArgs)
		} else {
			myTeamRouteSS(s, sender, "GET")
		}
	case "team":
		getTeamRouteSS(s, cmdArgs)
	case "pcprofile":
		profileRouteOW(s, cmdArgs, "pc")
	case "pcstats":
		statsRouteOW(s, cmdArgs, "pc")
	case "psprofile":
		profileRouteOW(s, cmdArgs, "psn")
	case "psstats":
		statsRouteOW(s, cmdArgs, "psn")
	case "xprofile":
		profileRouteOW(s, cmdArgs, "xbl")
	case "xstats":
		statsRouteOW(s, cmdArgs, "xbl")
	default:
		messageSend(s, "Enter a valid command")
	}
}
