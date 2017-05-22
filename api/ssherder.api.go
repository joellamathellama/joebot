package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "reflect"
	// "regexp"
	"joebot/rds"
	"joebot/tools"
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

type expectedStones struct {
	ID        int    `json:"id"`
	Element   string `json:"element"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	Rarity    int    `json:"rarity"`
	Type      string `json:"type"`
	Zodiac    string `json:"zodiac"`
	EvolvesTo int    `json:"evolves_to"`
	Skills    []int  `json:"skills"`
}

// HOST: https://ssherder.com
// Players: /data-api/characters/
// Skills: /data-api/skills/

func GetPlayers() {
	res, err := http.Get("https://ssherder.com/data-api/characters/")
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Somethings wrong with Ssherder!")
		return
	}
	defer res.Body.Close()

	// ReadAll to a byte array for Unmarshal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Error with: ioutil.ReadAll(res.Body)")
		return
	}

	// Unmarshal JSON data into struct
	var playerStruct []expectedPlayers
	if err := json.Unmarshal(body, &playerStruct); err != nil {
		tools.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(body, &playerStruct)")
		return
	}

	// loop and store
	for i := 0; i < len(playerStruct); i++ {
		playerMap := make(map[string]string)
		playerMap["Story"] = playerStruct[i].Story
		playerMap["ID"] = strconv.Itoa(playerStruct[i].ID)
		playerMap["Stones"] = strings.Join(playerStruct[i].Stones, ", ")

		var (
			skillString string
			ace         string
			active      string
			passives    string
		)

		for k := 0; k < len(playerStruct[i].Skills); k++ {
			// Define hash key, HGetAll, assign skill info
			hashKey := "skill_" + strconv.Itoa(playerStruct[i].Skills[k])

			skillHash, err := rds.RC.HGetAll(hashKey).Result()
			if err != nil {
				tools.WriteErr(err)
				fmt.Println("Error getting Skill Hash")
			}

			skillName := skillHash["Name"]
			skillDesc := skillHash["Description"]
			skillCat := skillHash["Category"]
			skillCost := skillHash["SpiritCost"]
			skillCD := skillHash["Cooldown"]

			// How I want one line printed
			skillInfo := fmt.Sprintf("**%s** [%s] \n%s\n\n", skillName, strings.ToLower(skillCat), skillDesc)

			if skillCat == "ace" {
				ace = skillInfo
			} else if skillCat == "active" { // active skills have a unique print
				active = fmt.Sprintf("**%s** [%s, %s spirit, %sm] \n%s\n\n", skillName, strings.ToLower(skillCat), skillCost, skillCD, skillDesc)
			} else { // Multiple passives per player
				passives = passives + skillInfo
			}
		}
		// Order I want it all in after ace: active > passives
		skillString = ace
		skillString = skillString + active
		skillString = skillString + passives

		playerMap["Skills"] = skillString

		// Example name: "Z101 Raklet"
		// Split it: ["Z101", "Raklet"]
		// Create the same player entries the keys: "Z101 Raklet", "Z101", and "Raklet"
		playerName := playerStruct[i].Name
		splitName := strings.Split(playerName, " ")

		stringID := strconv.Itoa(playerStruct[i].ID) // stringify ID
		keyID := string(stringID[0])                 // grab first index in string form
		lookupKey := playerStruct[i].Name + "_" + keyID

		// Store Character's name by Ssherder IDs
		rds.RedisSet(rds.RC, playerMap["ID"], playerName)

		// set full name, then loop over(if two or more) splitName
		rds.RC.HMSet(strings.ToLower(lookupKey), playerMap)
		if len(splitName) > 1 {
			for x := 0; x < len(splitName); x++ {
				// check if it exists already
				splitKey := fmt.Sprintf("%s_%s", strings.ToLower(splitName[x]), keyID)
				exists, err := rds.RC.Exists(splitKey).Result()
				if err != nil {
					tools.WriteErr(err)
					return
				} else if !exists {
					rds.RC.HMSet(splitKey, playerMap)
				} else {
					continue
				}
			}
		}
	}

	// Skill info is saved, now...
	// Also save other player info under a UUID key

	// _, err := io.Copy(os.Stdout, res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func GetSkills() {
	res, err := http.Get("https://ssherder.com/data-api/skills/")
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Api Server is down!")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Error with: ioutil.ReadAll(res.Body)")
		return
	}

	var skillsStruct []expectedSkills
	if err := json.Unmarshal(body, &skillsStruct); err != nil {
		tools.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(body, &skillsStruct)")
		return
	}

	// ID(stringified) to lookup
	for i := 0; i < len(skillsStruct); i++ {
		skillsMap := make(map[string]string)
		skillsMap["Name"] = skillsStruct[i].Name
		skillsMap["Category"] = skillsStruct[i].Category
		skillsMap["Cooldown"] = strconv.Itoa(skillsStruct[i].Cooldown)
		skillsMap["SpiritCost"] = skillsStruct[i].SpiritCost

		// skillsMap["Description"] = skillsStruct[i].Description
		editDesc := skillsStruct[i].Description

		// ch := make(chan bool)
		// finalDesc := replaceDesc(editDesc, skillsStruct[i].Effects, ch)
		finalDesc := replaceDesc(editDesc, skillsStruct[i].Effects, skillsStruct[i].Category)
		// <-ch
		// // fmt.Println(finalDesc)
		skillsMap["Description"] = finalDesc

		lookupKey := "skill_" + strconv.Itoa(skillsStruct[i].ID)
		rds.RC.HMSet(lookupKey, skillsMap)
	}
}

/*
	replaceDesc() replaces all variables in skill description string with correct values
	[#][1] + [#][2] * 4 = skill maxed out
*/
func replaceDesc(s string, i [][]interface{}, cat string) string {
	final := s
	for x := 0; x < len(i); x++ {
		find := fmt.Sprintf("{%d}", x)
		base := i[x][1].(string)
		multi := i[x][2].(string)

		// convert to int for calculation
		baseVerted, _ := strconv.ParseFloat(base, 64)
		multiVerted, _ := strconv.ParseFloat(multi, 64)
		baseInt := int(baseVerted)
		var multiInt int
		if cat == "item" { // multiplier is higher for UQ stones("items")
			multiInt = int(multiVerted * 16)
		} else if cat == "ace" {
			multiInt = int(multiVerted)
		} else {
			multiInt = int(multiVerted * 4)
		}
		// convert back to string for string replacement
		replacement := strconv.Itoa((baseInt + multiInt))
		// If Burst Ace
		if cat == "ace" && multi != "0.000" {
			replacement = fmt.Sprintf("%d%%/%s", baseInt, replacement)
		}
		final = strings.Replace(final, find, replacement, -1)

		// fmt.Println(s, find, replacement, final)
	}
	return final
}

func GetStones() {
	res, err := http.Get("https://ssherder.com/data-api/stones/")
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Api Server is down!")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		tools.WriteErr(err)
		fmt.Println("Error with: ioutil.ReadAll(res.Body)")
		return
	}

	var stoneStruct []expectedStones
	if err := json.Unmarshal(body, &stoneStruct); err != nil {
		tools.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(body, &stoneStruct)")
		return
	}

	// concat description string and set in redis
	/*
		<Name>
		<Element> Unique
		<Skill 1>
		<Skill 2>
	*/

	for x := 0; x < len(stoneStruct); x++ {
		var (
			stoneName string
			skillDesc string
		)
		// if rarity 4+
		if stoneStruct[x].Rarity >= 4 {
			stoneName = stoneStruct[x].Name
			// fmt.Println(stoneName)
			skillDesc = fmt.Sprintf("**%s**\n", stoneName)
			skillDesc = fmt.Sprintf("%s%s *%s*\n", skillDesc, stoneStruct[x].Element, stoneStruct[x].Type)

			// HGET skill descriptions
			// Pretty up and store
			if len(stoneStruct[x].Skills) > 0 {
				for z := 0; z < len(stoneStruct[x].Skills); z++ {
					stoneSkillKey := fmt.Sprintf("skill_%d", stoneStruct[x].Skills[z])
					stoneSkillHash, err := rds.RC.HGetAll(stoneSkillKey).Result()
					if err != nil {
						tools.WriteErr(err)
					}
					skillDesc = fmt.Sprintf("%s%s\n", skillDesc, stoneSkillHash["Description"])
				}
			}
		}
		// Key example: stone_stone name, stone_stone, and stone_name
		// split the name if it contains any spaces
		if len(stoneName) != 0 {
			// stoneKey := fmt.Sprintf("stone_%s", strings.ToLower(stoneName))
			stoneKey := fmt.Sprintf("stone_%s", stoneName)
			splitName := strings.Split(stoneName, " ")
			ok := rds.RedisSet(rds.RC, stoneKey, skillDesc)
			if !ok {
				tools.WriteLog("Error: getStones() redisSet failed!")
			}
			if len(splitName) > 1 {
				for k := 0; k < len(splitName); k++ {
					splitKey := fmt.Sprintf("stone_%s", splitName[k])
					ok = rds.RedisSet(rds.RC, splitKey, strings.ToLower(skillDesc))
					if !ok {
						tools.WriteLog("Error: getStones() redisSet failed!")
					}
				}
			}
		}
	}
}

// func getChains() {
// 	res, err := http.Get("https://ssherder.com/data-api/chains/")
// 	if err != nil {
// 		tools.WriteErr(err)
// 		fmt.Println("Api Server is down!")
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		tools.WriteErr(err)
// 		fmt.Println("Error with: ioutil.ReadAll(res.Body)")
// 		return
// 	}

// }
