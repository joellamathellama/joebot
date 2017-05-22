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
	cmdResList["redditss"] = "http://reddit.com/r/soccerspirits"
	cmdResList["redditdwu"] = "http://reddit.com/r/dwunleashed"
	cmdResList["teamwork"] = "https://docs.google.com/spreadsheets/d/1x0Q4vUk_V3wUwzM5XR_66xytSbapoSFm_cHR9PYIERs/htmlview?sle=true#"
	cmdResList["chains"] = "https://ssherder.com/characters/#"
	cmdResList["help"] = "Shoutout to ssherder.com and api.lootbox.eu/documentation#/ for their APIs.\n\nTo talk to the bot, preface your message with '~jb'\n\n*General Commands:*\n**Write my own Note:** 'mynote <Your note here.>' (Ex. ~jb note sharr/elaine/renee...)\n**Read others Note:** 'note <Discord Name>' (Ex. ~jb note joellama)\n\n*Overwatch Commands:*\n**Lookup PC Profile:** 'PCprofile <Battlenet Tag>' (Ex. ~jb pcprofile joellama#1114)\n**Lookup PC Stats:** 'PCstats <Battlenet Tag>' (Ex. ~jb pcstats joellama#1114)\n**Lookup PS:** Same thing, except 'PSprofile, PSstats'\n**Lookup Xbox:** Same thing, except 'Xprofile, Xstats'\n\n*Soccer Spirits Commands:*\n**Lookup player info:** 'sstory, sstone, sslots, ssherder or sskills <Player Name>' (Ex. ~jb stats Griffith)\n**Quick links:** 'ourteams', 'apoc', 'reddit' (Ex. ~jb apoc)\n\n*Dynasty Warriors Unleashed Commands:*\n**Lookup Officer Legendary Passives:** 'dwup <Officer Name>' (Ex. ~jb dwup cai wenji)\n**Lookup Officer Stats:** 'dwus <Officer Name>' (Ex. ~jb dwus cai wenji)\n\n*Everything is case *insensitive!*(Except Bnet Tags)"
}

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageRoutes(s *dg.Session, m *dg.MessageCreate) {
	// Contents
	c := m.Content // full message sent by user
	nn := 3        //  bot command length(~jb)
	// Ignore all messages created by the bot itself and anything short of "~jb "
	if m.Author.ID == BotID {
		return
	} else if len(c) <= nn || tools.RegexpMatch("^(?i)(~JB)", c[0:nn]) != true {
		// fmt.Println("Not talking to bot")
		return
	}

	// split message by command and arguments
	cSplit := strings.Split(c, " ") // ["~jb", "command", [...]]
	cc := cSplit[1]
	cl := len(cc) + nn
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
	case "sstory":
		storyRouteSS(s, cmdArgs)
	case "sslots":
		slotesRouteSS(s, cmdArgs)
	case "ssherder":
		ssherderRouteSS(s, cmdArgs)
	case "sskills":
		skillsRouteSS(s, cmdArgs)
	case "sstone":
		stoneRouteSS(s, cmdArgs)
	case "mynote":
		if len(cmdArgs) > 0 {
			myTeamRouteSS(s, sender, cmdArgs)
		} else {
			myTeamRouteSS(s, sender, "GET")
		}
	case "note":
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
	case "dwup":
		passiveRouteDWU(s, cmdArgs)
	case "dwus":
		officerRouteDWU(s, cmdArgs)
	default:
		// messageSend(s, "Enter a valid command")
	}
}
