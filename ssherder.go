package main

import (

)

// Find a better name for this...
type expectedResData struct {
	ID int `json:"id"`
	ChainMissions []struct {
		Effect string `json:"effect"`
		Target int `json:"target"`
		BaseCharacters []int `json:"base_characters"`
	} `json:"chain_missions"`
}

func ssherderApi() {
	
}
