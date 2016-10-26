package main

import (
	// "fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	// "reflect"
	// "io"
	// "os"
)

// Expected JSON from Ssherder API
// It will come back in an Array of Objects
type expectedPlayers struct {
	ID              int      `json:"id"`
	ImageID         int      `json:"image_id"`
	BaseCharacter   int      `json:"base_character"`
	Name            string   `json:"name"`
	Cost            int      `json:"cost"`
	Element         string   `json:"element"`
	Gender          string   `json:"gender"`
	Rarity          int      `json:"rarity"`
	Category        string   `json:"category"`
	Role            string   `json:"role"`
	Season          int      `json:"season"`
	Stones          []string `json:"stones"`
	MinPow          int      `json:"min_pow"`
	MinTec          int      `json:"min_tec"`
	MinVit          int      `json:"min_vit"`
	MinSpd          int      `json:"min_spd"`
	MaxPow          int      `json:"max_pow"`
	MaxTec          int      `json:"max_tec"`
	MaxVit          int      `jdson:"max_vit"`
	MaxSpd          int      `json:"max_spd"`
	Story           string   `json:"story"`
	WeatherImmunity string   `json:"weather_immunity"`
	Illustrator     int      `json:"illustrator"`
	VoiceActor      int      `json:"voice_actor"`
	IsLegend        bool     `json:"is_legend"`
	IsSpecial       bool     `json:"is_special"`
	Skills          []int    `json:"skills"`
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
	for i := 0; i < len(playerStruct); i++ {
		playerMap := make(map[string]string)
		playerMap["Story"] = playerStruct[i].Story
		playerMap["ID"] = strconv.Itoa(playerStruct[i].ID)
		playerMap["Stones"] = strings.Join(playerStruct[i].Stones, ", ")

		// instead of storing skill id's([#, #, #, #, #])
		// do a lookup on those ids to get the skill data
		// concat all the skill data, and then set it in the playerMap
		// use callback to make sure skills are called first, then this player call
		skillString := ""
		ace := ""
		active := ""
		passives := ""
		for k := 0; k < len(playerStruct[i].Skills); k++ {
			// TODO: Use HGetAll instead

			nameKey := "skill_" + strconv.Itoa(playerStruct[i].Skills[k])
			skillName, err := rc.HGet(nameKey, "Name").Result()
			if err != nil {
				panic(err)
			} else {
				// fmt.Println(skillName)
			}

			descKey := "skill_" + strconv.Itoa(playerStruct[i].Skills[k])
			skillDesc, err := rc.HGet(descKey, "Description").Result()
			if err != nil {
				panic(err)
			} else {
				// fmt.Println(skillDesc)
			}

			catKey := "skill_" + strconv.Itoa(playerStruct[i].Skills[k])
			skillCat, err := rc.HGet(catKey, "Category").Result()
			if err != nil {
				panic(err)
			} else {
				// fmt.Println(skillCat)
			}

			skillInfo := skillName + ": " + skillDesc + "\n\n"

			if skillCat == "ace" { // Only one ace and active per player
				ace = skillInfo
			} else if skillCat == "active" {
				active = skillInfo
			} else { // Three passive skills per player
				passives = passives + skillInfo
			}

			skillString = ace
			skillString = skillString + active
			skillString = skillString + passives
		}
		playerMap["Skills"] = skillString

		stringID := strconv.Itoa(playerStruct[i].ID) // stringify ID
		keyID := string(stringID[0])                 // grab first index in string form
		lookupKey := playerStruct[i].Name + "_" + keyID

		rc.HMSet(lookupKey, playerMap)
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

	var skillsStruct []expectedSkills
	json.Unmarshal(body, &skillsStruct)

	// ID(stringified) to lookup
	for i := 0; i < len(skillsStruct); i++ {
		skillsMap := make(map[string]string)
		skillsMap["Name"] = skillsStruct[i].Name
		skillsMap["Description"] = skillsStruct[i].Description
		// skillsMap["Effects"] = skillsStruct[i].Effects
		skillsMap["Category"] = skillsStruct[i].Category

		// fmt.Println(skillsStruct[i].Name, reflect.TypeOf(skillsStruct[i].Name))
		// fmt.Println(skillsStruct[i].Description, reflect.TypeOf(skillsStruct[i].Description))
		// fmt.Println(skillsStruct[i].Category, reflect.TypeOf(skillsStruct[i].Category))

		lookupKey := "skill_" + strconv.Itoa(skillsStruct[i].ID)
		// fmt.Println(lookupKey, reflect.TypeOf(lookupKey))
		rc.HMSet(lookupKey, skillsMap)
	}
}
