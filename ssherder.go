package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "reflect"
	// "regexp"
	"strconv"
	"strings"
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

			// How I want it printed
			skillInfo := fmt.Sprintf("**%s** [%s] \n%s\n\n", skillName, skillCat, skillDesc)

			// Order I want it printed
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
		skillsMap["Category"] = skillsStruct[i].Category

		// skillsMap["Description"] = skillsStruct[i].Description
		editDesc := skillsStruct[i].Description

		// [#][1] + [#][2] * 4 = skill maxed out
		// recursively look for {#} and replace it based on the #
		// ch := make(chan bool)
		// finalDesc := replaceDesc(editDesc, skillsStruct[i].Effects, ch)
		finalDesc := replaceDesc(editDesc, skillsStruct[i].Effects)
		// <-ch
		// // fmt.Println(finalDesc)
		skillsMap["Description"] = finalDesc

		lookupKey := "skill_" + strconv.Itoa(skillsStruct[i].ID)
		// fmt.Println(lookupKey, reflect.TypeOf(lookupKey))
		rc.HMSet(lookupKey, skillsMap)
	}
}

func replaceDesc(s string, i [][]interface{}) string {
	final := s
	for x := 0; x < len(i); x++ {
		find := fmt.Sprintf("{%d}", x)
		base := i[x][1].(string)
		multi := i[x][2].(string)

		// fmt.Println(x, find, len(i), base, multi)

		// convert to int for calculation
		baseVerted, _ := strconv.ParseFloat(base, 64)
		multiVerted, _ := strconv.ParseFloat(multi, 64)
		baseInt := int(baseVerted)
		multiInt := int(multiVerted * 4)
		// calculate then convert back to string for string replacement
		replacement := strconv.Itoa((baseInt + multiInt))
		final = strings.Replace(final, find, replacement, -1)

		// fmt.Println(s, find, replacement, final)
	}
	return final
}

// FAILED ATTEMPT AT USING REGEX + RECURSION :C
// Keeping it cause why not
// func replaceDesc(s string, i [][]interface{}) string {
// 	re := regexp.MustCompile("[{][0-9][}]")
// 	a := re.FindStringIndex(s) // [start_index end_index]
// 	// if no index exit
// 	if len(a) == 0 {
// 		return s
// 	}
// 	// grab the # between the {}
// 	b := a[0] + 1
// 	// type assertion to access the stringified numbers
// 	base := ""
// 	multi := ""
// 	f, _ := strconv.Atoi(string(s[b]))

// 	if len(i) >= (f + 1) {
// 		base = i[f][1].(string)
// 		multi = i[f][2].(string)
// 	} else {
// 		base = "0"
// 		multi = "0"
// 	}

// 	// convert to int for calculation
// 	baseVerted, _ := strconv.ParseFloat(base, 64)
// 	multiVerted, _ := strconv.ParseFloat(multi, 64)
// 	baseInt := int(baseVerted)
// 	multiInt := int(multiVerted * 4)

// 	// calculate then convert back to string for string replacement
// 	finalNum := strconv.Itoa((baseInt + multiInt))

// 	c := re.ReplaceAllString(s, finalNum)
// 	return replaceDesc(c, i)
// }
