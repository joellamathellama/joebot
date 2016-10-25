package main

import (
	// "fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	// "reflect"
	// "io"
	// "os"
)

// Expected JSON from Ssherder API
// It will come back in an Array of Objects
type expectedPlayers struct {
	ID int `json:"id"`
	ImageID int `json:"image_id"`
	BaseCharacter int `json:"base_character"`
	Name string `json:"name"`
	Cost int `json:"cost"`
	Element string `json:"element"`
	Gender string `json:"gender"`
	Rarity int `json:"rarity"`
	Category string `json:"category"`
	Role string `json:"role"`
	Season int `json:"season"`
	Stones []string `json:"stones"`
	MinPow int `json:"min_pow"`
	MinTec int `json:"min_tec"`
	MinVit int `json:"min_vit"`
	MinSpd int `json:"min_spd"`
	MaxPow int `json:"max_pow"`
	MaxTec int `json:"max_tec"`
	MaxVit int `json:"max_vit"`
	MaxSpd int `json:"max_spd"`
	Story string `json:"story"`
	WeatherImmunity string `json:"weather_immunity"`
	Illustrator int `json:"illustrator"`
	VoiceActor int `json:"voice_actor"`
	IsLegend bool `json:"is_legend"`
	IsSpecial bool `json:"is_special"`
	Skills []int `json:"skills"`
}

type expectedSkills struct {
	Accumable       bool            `json:"accumable"`
	AccumableTo     interface{}     `json:"accumable_to"`
	Category        string          `json:"category"`
	Cooldown        int             `json:"cooldown"`
	CooldownGrowth  int             `json:"cooldown_growth"`
	CooldownInitial int             `json:"cooldown_initial"`
	Description     string          `json:"description"`
	Duration        int             `json:"duration"`
	Effects         [][]interface{} `json:"effects"`
	Icon            string          `json:"icon"`
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	SpiritCost      string          `json:"spirit_cost"`
}

// HOST: https://ssherder.com
// Players: /data-api/characters/
// Skills: /data-api/skills/

func getPlayers() {
	res, err := http.Get("https://ssherder.com/data-api/characters/")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// ReadAll to a byte array for Unmarshal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
	    panic(err.Error())
	}

	// Unmarshal JSON data into struct
	var playerStruct []expectedPlayers
	json.Unmarshal(body, &playerStruct)

	// loop and store
	for i := 0; i < len(playerStruct); i++{
		playerMap := make(map[string]string)
		playerMap["Story"] = playerStruct[i].Story
		playerMap["Stones"] = strings.Join(playerStruct[i].Stones, ", ")

		rc.HMSet(playerStruct[i].Name, playerMap)

		for x := 0; x < len(playerStruct[i].Skills); x++{
			playerName := playerStruct[i].Name
			createdKey := playerName + "_skills"

			// key = playername_skills, value = list of stringified(int -> string) skills
			rc.LPush(createdKey, strconv.Itoa(playerStruct[i].Skills[x]))
		}
	}

	// _, err := io.Copy(os.Stdout, res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func getSkills() {
	res, err := http.Get("https://ssherder.com/data-api/skills/")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
	    panic(err.Error())
	}

	var skillStruct []expectedSkills
	json.Unmarshal(body, &skillStruct)

	// fmt.Println(rc.HGetAll(strconv.Itoa(skillStruct[0].ID)))

	// ID(stringified) to lookup
	// Store Name, Description
	for i := 0; i < len(skillStruct); i++{
		skillMap := make(map[string]string)
		skillMap["Name"] = skillStruct[i].Name
		skillMap["Description"] = skillStruct[i].Description

		rc.HMSet(strconv.Itoa(skillStruct[i].ID), skillMap)
	}
}
