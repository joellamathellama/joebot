package main

import (

)

type expectedData struct {
	ID int `json:"id"`
	ChainMissions []struct {
		Effect string `json:"effect"`
		Target int `json:"target"`
		BaseCharacters []int `json:"base_characters"`
	} `json:"chain_missions"`
}

func ssherderApi() {
	
}
