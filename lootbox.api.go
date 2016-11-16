package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "reflect"
)

// type expectedPlayerProfile struct {
// 	Data struct {
// 		Username string `json:"username"`
// 		Level    int    `json:"level"`
// 		Games    struct {
// 			Quick struct {
// 				Wins string `json:"wins"`
// 			} `json:"quick"`
// 			Competitive struct {
// 			} `json:"competitive"`
// 		} `json:"games"`
// 		Playtime struct {
// 			Quick string `json:"quick"`
// 		} `json:"playtime"`
// 		Avatar      string `json:"avatar"`
// 		Competitive struct {
// 			Rank interface{} `json:"rank"`
// 		} `json:"competitive"`
// 		LevelFrame string `json:"levelFrame"`
// 		Star       string `json:"star"`
// 	} `json:"data"`
// }

// type expectedPlayerStats struct {
// 	MeleeFinalBlows            string `json:"MeleeFinalBlows"`
// 	SoloKills                  string `json:"SoloKills"`
// 	ObjectiveKills             string `json:"ObjectiveKills"`
// 	FinalBlows                 string `json:"FinalBlows"`
// 	DamageDone                 string `json:"DamageDone"`
// 	Eliminations               string `json:"Eliminations"`
// 	EnvironmentalKill          string `json:"EnvironmentalKill"`
// 	Multikills                 string `json:"Multikills"`
// 	HealingDone                string `json:"HealingDone"`
// 	ReconAssists               string `json:"ReconAssists"`
// 	EliminationsMostinGame     string `json:"Eliminations-MostinGame"`
// 	FinalBlowsMostinGame       string `json:"FinalBlows-MostinGame"`
// 	DamageDoneMostinGame       string `json:"DamageDone-MostinGame"`
// 	HealingDoneMostinGame      string `json:"HealingDone-MostinGame"`
// 	DefensiveAssistsMostinGame string `json:"DefensiveAssists-MostinGame"`
// 	OffensiveAssistsMostinGame string `json:"OffensiveAssists-MostinGame"`
// 	ObjectiveKillsMostinGame   string `json:"ObjectiveKills-MostinGame"`
// 	ObjectiveTimeMostinGame    string `json:"ObjectiveTime-MostinGame"`
// 	MultikillBest              string `json:"Multikill-Best"`
// 	SoloKillsMostinGame        string `json:"SoloKills-MostinGame"`
// 	TimeSpentonFireMostinGame  string `json:"TimeSpentonFire-MostinGame"`
// 	MeleeFinalBlowsAverage     string `json:"MeleeFinalBlows-Average"`
// 	TimeSpentonFireAverage     string `json:"TimeSpentonFire-Average"`
// 	SoloKillsAverage           string `json:"SoloKills-Average"`
// 	ObjectiveTimeAverage       string `json:"ObjectiveTime-Average"`
// 	ObjectiveKillsAverage      string `json:"ObjectiveKills-Average"`
// 	HealingDoneAverage         string `json:"HealingDone-Average"`
// 	FinalBlowsAverage          string `json:"FinalBlows-Average"`
// 	DeathsAverage              string `json:"Deaths-Average"`
// 	DamageDoneAverage          string `json:"DamageDone-Average"`
// 	EliminationsAverage        string `json:"Eliminations-Average"`
// 	Deaths                     string `json:"Deaths"`
// 	EnvironmentalDeaths        string `json:"EnvironmentalDeaths"`
// 	Cards                      string `json:"Cards"`
// 	Medals                     string `json:"Medals"`
// 	MedalsGold                 string `json:"Medals-Gold"`
// 	MedalsSilver               string `json:"Medals-Silver"`
// 	MedalsBronze               string `json:"Medals-Bronze"`
// 	GamesWon                   string `json:"GamesWon"`
// 	TimeSpentonFire            string `json:"TimeSpentonFire"`
// 	ObjectiveTime              string `json:"ObjectiveTime"`
// 	TimePlayed                 string `json:"TimePlayed"`
// 	MeleeFinalBlowMostinGame   string `json:"MeleeFinalBlow-MostinGame"`
// 	DefensiveAssists           string `json:"DefensiveAssists"`
// 	DefensiveAssistsAverage    string `json:"DefensiveAssists-Average"`
// 	OffensiveAssists           string `json:"OffensiveAssists"`
// 	OffensiveAssistsAverage    string `json:"OffensiveAssists-Average"`
// 	ReconAssistsAverage        string `json:"ReconAssists-Average"`
// 	ReconAssistMostinGame      string `json:"ReconAssist-MostinGame"`
// }

// Example Player info API: https://api.lootbox.eu/pc/us/jawnkeem-1982/profile

func getPlayerProfile(bnetID string, platform string) string {
	url := fmt.Sprintf("https://api.lootbox.eu/%s/us/%s/profile", platform, bnetID)
	res, err := http.Get(url)
	if err != nil {
		writeErr(err)
		return "Api Server is down!"
	}
	defer res.Body.Close()

	// ReadAll to a byte array for Unmarshal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		writeErr(err)
		return fmt.Sprintf("**error:** Found no user with the BattleTag: %s", strings.Replace(bnetID, "-", "#", -1))
	}

	// Unmarshal JSON data into struct
	// var profileStruct []expectedPlayerProfile
	var profileStruct interface{}
	if err := json.Unmarshal(body, &profileStruct); err != nil {
		writeErr(err)
		return fmt.Sprintf("**error:** Found no user with the BattleTag: %s", strings.Replace(bnetID, "-", "#", -1))
	}

	m := profileStruct.(map[string]interface{})
	if m["data"] == nil {
		return fmt.Sprintf("**error:** Found no user with the BattleTag: %s", strings.Replace(bnetID, "-", "#", -1))
	}

	// Type assertion
	n := m["data"].(map[string]interface{})
	level := n["level"].(float64)
	games := n["games"].(map[string]interface{})
	quicks := games["quick"].(map[string]interface{})
	playtime := n["playtime"].(map[string]interface{})
	// Not everyone plays competitive
	comprank := n["competitive"].(map[string]interface{})
	compranking := comprank["rank"]
	if compranking == nil {
		compranking = "N/A"
	}
	comp := games["competitive"].(map[string]interface{})
	compwins := comp["wins"]
	if compwins == nil {
		compwins = "0"
	}
	compplayed := comp["played"]
	if compplayed == nil {
		compplayed = "0"
	}
	comptime := playtime["competitive"]
	if comptime == nil {
		comptime = "0"
	}

	messageToSend := ""

	messageToSend = fmt.Sprintf("**Username:** %s\n", n["username"])
	messageToSend = fmt.Sprintf("%s**Level:** %d\n", messageToSend, int(level))
	messageToSend = fmt.Sprintf("%s**Quick Wins:** %s\n", messageToSend, quicks["wins"])
	messageToSend = fmt.Sprintf("%s**Quick Time:** %s\n", messageToSend, playtime["quick"])
	messageToSend = fmt.Sprintf("%s**Competitive Rank:** %s\n", messageToSend, compranking)
	messageToSend = fmt.Sprintf("%s**Competitive Wins:** %s\n", messageToSend, compwins)
	messageToSend = fmt.Sprintf("%s**Competitive Played:** %s\n", messageToSend, compplayed)
	messageToSend = fmt.Sprintf("%s**Competitive Time:** %s\n", messageToSend, comptime)

	// Save into redis
	playerHash := fmt.Sprintf("%s%s", bnetID, platform)
	redisMap := make(map[string]string)
	redisMap["profile"] = messageToSend
	rc.HMSet(playerHash, redisMap)

	return messageToSend
}

func getPlayerStats(bnetID string, platform string) string {
	url := fmt.Sprintf("https://api.lootbox.eu/%s/us/%s/quick-play/allHeroes/", platform, bnetID)
	res, err := http.Get(url)
	if err != nil {
		writeErr(err)
		return "Somethings wrong with lootbox!"
	}
	defer res.Body.Close()

	// ReadAll to a byte array for Unmarshal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		writeErr(err)
		return fmt.Sprintf("**error:** Found no user with the BattleTag: %s", strings.Replace(bnetID, "-", "#", -1))
	}

	// Unmarshal JSON data into struct
	// var statsStruct []expectedPlayerStats
	// if err := json.Unmarshal(body, &statsStruct); err == nil {
	// 	panic(err)
	// }

	var statsStruct interface{}
	if err := json.Unmarshal(body, &statsStruct); err != nil {
		writeErr(err)
		return fmt.Sprintf("**error:** Found no user with the BattleTag: %s", strings.Replace(bnetID, "-", "#", -1))
	}

	messageToSend := ""

	m := statsStruct.(map[string]interface{})

	for key, value := range m {
		if s, ok := value.(string); ok {
			// fmt.Printf("%q is a string: %q\n", key, s)
			messageToSend = fmt.Sprintf("%s**%s:** %s\n", messageToSend, key, s)
		}
	}

	playerHash := fmt.Sprintf("%s%s", bnetID, platform)
	redisMap := make(map[string]string)
	redisMap["stats"] = messageToSend
	rc.HMSet(playerHash, redisMap)

	return messageToSend
}
