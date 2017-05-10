package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"joebot/rds"
	t "joebot/tools"
	"strings"
)

type ExpectedOfficer struct {
	OfficerName string      `json:"Officer name"`
	Type        string      `json:"Type"`
	Element     string      `json:"Element"`
	BaseStar    json.Number `json:"Base Star"`
	HP          json.Number `json:"HP"`
	ATK         json.Number `json:"ATK"`
	DEF         json.Number `json:"DEF"`
	CritChance  json.Number `json:"Crit Chance"`
	CritDodge   json.Number `json:"Crit Dodge"`
	Farmable    string      `json:"Farmable"`
	Skill1      string      `json:"Skill 1"`
	Skill2      string      `json:"Skill 2"`
}

func OfficersToRedis() {
	filePath := "./json/dwu/officers.stats.json"
	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.WriteErr(err)
		eMess := fmt.Sprintf("Somethings wrong with reading %s!", filePath)
		fmt.Println(eMess)
		return
	}

	var oS []ExpectedOfficer
	if err := json.Unmarshal(res, &oS); err != nil {
		t.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(res, &officerStruct)")
		return
	} else {
		t.WriteLog("Successful: officers.stats.json Unmarshal to officerStruct")
	}

	for i := 0; i < len(oS); i++ {
		officerKey := fmt.Sprintf("officer_%s", strings.ToLower(oS[i].OfficerName))
		statString := fmt.Sprintf("**%s** [Natural %s]\n%s / %s\n__Max Lvl Stats__\nHP: %s\nATK: %s\nDEF: %s\nCrit Chance: %s\nCrit Dodge: %s\n__Leader Skills__\nBase: %s\n6 Star: %s\n", oS[i].OfficerName, oS[i].BaseStar, oS[i].Type, oS[i].Element, oS[i].HP, oS[i].ATK, oS[i].DEF, oS[i].CritChance, oS[i].CritDodge, oS[i].Skill1, oS[i].Skill2)

		rLog := fmt.Sprintf("RedisSet(%s,...)", officerKey)
		t.WriteLog(rLog)
		rds.RedisSet(rds.RC, officerKey, statString)
	}
}

func PassivesToRedis(fileName string) {
	wLog := fmt.Sprintf("PassivesToRedis(%s)", fileName)
	t.WriteLog(wLog)
	filePath := fmt.Sprintf("./json/dwu/%s.passives.json", fileName)
	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.WriteErr(err)
		eMess := fmt.Sprintf("Somethings wrong with reading %s!", filePath)
		fmt.Println(eMess)
		return
	}

	var PassiveStruct []interface{}
	if err := json.Unmarshal(res, &PassiveStruct); err != nil {
		t.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(body, &PassiveStruct)")
		return
	} else {
		t.WriteLog("Successful: passives.json Unmarshal to PassiveStruct")
	}

	officerMap := make(map[string]string)

	for i := 0; i < len(PassiveStruct); i++ {
		var (
			on string
			pn string
			sd string
			pl string
			ss string
		)

		m := PassiveStruct[i].(map[string]interface{})
		for k, v := range m {
			val := v.(string)
			if k == "Officer Name" {
				val = strings.ToLower(val)
				on = fmt.Sprintf("passive_%s", val)
			} else if k == "Passive Name" {
				pn = val
			} else if k == "Skill Description" {
				sd = val
			} else if k == "Passive Level" {
				pl = val
			}
		}

		ss = fmt.Sprintf("**%s** [Lvl %s]\n%s.\n\n", pn, pl, sd)

		if val, ok := officerMap[on]; ok {
			officerMap[on] = fmt.Sprintf("%s%s", val, ss)
		} else {
			officerMap[on] = ss
		}
	}

	for k, v := range officerMap {
		rLog := fmt.Sprintf("RedisSet(%s,...)", k)
		t.WriteLog(rLog)
		rds.RedisSet(rds.RC, k, v)
	}
}
