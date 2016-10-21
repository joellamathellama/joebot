package main

import (
	// "fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	// "reflect"
	// "io"
	// "os"
)

// Expected JSON from Ssherder API
// It will come back in an Array of Objects
type expectedChar struct {
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

// HOST: https://ssherder.com
// Characters: /data-api/characters/

func getChars() {
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
	var createdStruct []expectedChar
	json.Unmarshal(body, &createdStruct)

	// loop and store
	for i := 0; i < len(createdStruct); i++{
		myMap := make(map[string]string)
		myMap["Story"] = createdStruct[i].Story
		myMap["Stones"] = strings.Join(createdStruct[i].Stones, ", ")

		redisClient.HMSet(createdStruct[i].Name, myMap)
	}

	// _, err := io.Copy(os.Stdout, res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
