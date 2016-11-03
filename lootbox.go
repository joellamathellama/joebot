package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "strings"
	// "reflect"
)

// type expectedPlayerStats struct {
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

// Example Player info API: https://api.lootbox.eu/pc/us/jawnkeem-1982/profile

func getPlayerStats(bnetID string) string {
	url := fmt.Sprintf("https://api.lootbox.eu/pc/us/%s/profile", bnetID)
	res, err := http.Get(url)
	if err != nil {
		return "Not Found."
	}
	defer res.Body.Close()

	// ReadAll to a byte array for Unmarshal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "Not Found."
	}

	// Unmarshal JSON data into struct
	// var playerStruct []expectedPlayerStats
	var playerStruct interface{}
	if err := json.Unmarshal(body, &playerStruct); err != nil {
		return "Not Found."
	}

	m := playerStruct.(map[string]interface{})
	n := m["data"].(map[string]interface{})
	level := n["level"].(float64)
	games := n["games"].(map[string]interface{})
	quicks := games["quick"].(map[string]interface{})
	playtime := n["playtime"].(map[string]interface{})
	comp := games["competitive"].(map[string]interface{})
	compwins := comp["wins"]
	if compwins == nil {
		compwins = "0"
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
	messageToSend = fmt.Sprintf("%s**Competitive Wins:** %s\n", messageToSend, compwins)
	messageToSend = fmt.Sprintf("%s**Competitive Time:** %s\n", messageToSend, comptime)

	return messageToSend
}
