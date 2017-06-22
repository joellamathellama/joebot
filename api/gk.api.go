package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"joebot/rds"
	t "joebot/tools"
	"log"
	"os"
	"strings"
)

type ExpectedPilot struct {
	ID                  json.Number `json:"id"`
	SIdx                json.Number `json:"S_Idx"`
	CType               json.Number `json:"C_Type"`
	UnitID              json.Number `json:"unitId"`
	ResourceID          string      `json:"resourceId"`
	IsAcademyExposure   string      `json:"_isAcademyExposure"`
	AcademyExposureRate json.Number `json:"academyExposureRate"`
	Grade               json.Number `json:"grade"`
	RecruitHonor        json.Number `json:"recruitHonor"`
	RecruitGold         json.Number `json:"recruitGold"`
	OverlapReward       json.Number `json:"overlapReward"`
	MedalExplanation    json.Number `json:"medalExplanation"`
	Favormax            json.Number `json:"favormax"`
	Explanation         json.Number `json:"explanation"`
}

func LocalizationToRedis() {
	f, _ := os.Open("./json/Regulation/Localization.txt")
	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error in parsing :", err)
			}
			break
		}
		sLine := strings.Split(line, "|")
		i := sLine[0]
		en := sLine[2]
		rKey := fmt.Sprintf("gk_en_%s", i)
		rds.RedisSet(rds.RC, rKey, en)
	}
}

func PilotsToRedis() {
	filePath := "./json/Regulation/CommanderDataTable.json"
	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.WriteErr(err)
		eMess := fmt.Sprintf("Somethings wrong with reading %s!", filePath)
		fmt.Println(eMess)
		return
	}

	var eP []ExpectedPilot
	if err := json.Unmarshal(res, &eP); err != nil {
		t.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(res, &ExpectedPilot)")
		return
	} else {
		t.WriteLog("Successful: Regulation/CommanderDataTable.json Unmarshal to ExpectedPilot")
	}

	// AVOID SIdx: 1008021(Empire Commander), 1008019(Empire Soldier)
	for i := 0; i < len(eP); i++ {
		if eP[i].SIdx != "1008021" && eP[i].SIdx != "1008019" {
			pName, _ := rds.RedisGet(rds.RC, fmt.Sprintf("gk_en_%s", eP[i].SIdx))
			pNameSplit := strings.Split(pName, " ")
			pSkillDesc := ""
			// each pilot has 4 skills(append redis key with 1-4 accordingly)
			for k := 1; k <= 4; k++ {
				sName, _ := rds.RedisGet(rds.RC, fmt.Sprintf("gk_en_10%s%d", eP[i].UnitID, k))
				sDesc, _ := rds.RedisGet(rds.RC, fmt.Sprintf("gk_en_20%s%d", eP[i].UnitID, k))
				sDesc = strings.Replace(sDesc, `\n\n`, "\n", -1)
				sDesc = strings.Replace(sDesc, `\n`, "\n", -1)
				pSkillDesc = fmt.Sprintf("%s**%s**\n%s\n\n", pSkillDesc, sName, sDesc)
			}

			if len(pSkillDesc) > 30 {
				pLow := strings.ToLower(pName)
				rds.RedisSet(rds.RC, fmt.Sprintf("gk_ps_%s", pLow), pSkillDesc)
				// Also set first and last individually
				for j := 0; j < len(pNameSplit); j++ {
					pL := strings.ToLower(pNameSplit[j])
					rds.RedisSet(rds.RC, fmt.Sprintf("gk_ps_%s", pL), pSkillDesc)
				}
			}
		}
	}
}
