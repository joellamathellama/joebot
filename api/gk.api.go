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

type ExpectedGKSkill struct {
	Key                 json.Number   `json:"key"`
	OpenGrade           json.Number   `json:"openGrade"`
	OrderType           json.Number   `json:"orderType"`
	TargetableRange     string        `json:"targetableRange"`
	SkillName           json.Number   `json:"skillName"`
	SkillDescription    json.Number   `json:"skillDescription"`
	Thumbnail           string        `json:"thumbnail"`
	Rangetype           json.Number   `json:"rangetype"`
	Rangeicon           string        `json:"rangeicon"`
	CutInEffectID       json.Number   `json:"cutInEffectId"`
	UnitMotionDrk       string        `json:"unitMotionDrk"`
	FirePatterns        []string      `json:"firePatterns"`
	ProjectileDrks      []json.Number `json:"projectileDrks"`
	CooldownPoint       json.Number   `json:"cooldownPoint"`
	Accuracy            json.Number   `json:"accuracy"`
	AttackDamage        json.Number   `json:"attackDamage"`
	CriticalChance      json.Number   `json:"criticalChance"`
	CriticalChanceBonus json.Number   `json:"criticalChanceBonus"`
	Depletion           json.Number   `json:"depletion"`
	Healing             json.Number   `json:"healing"`
	UnitLife            json.Number   `json:"unitLife"`
	RemainedHealthRate  json.Number   `json:"remainedHealthRate"`
	InitSp              json.Number   `json:"initSp"`
	MaxSp               json.Number   `json:"maxSp"`
	SpCostOnEnter       json.Number   `json:"spCostOnEnter"`
	SpCostOnTurn        json.Number   `json:"spCostOnTurn"`
	SpCostOnBeHit       json.Number   `json:"spCostOnBeHit"`
	SpCostOnCombo       json.Number   `json:"spCostOnCombo"`
	SpOnHit             json.Number   `json:"spOnHit"`
	SpOnCriticalHit     json.Number   `json:"spOnCriticalHit"`
	SpOnBeHit           json.Number   `json:"spOnBeHit"`
	SpOnDestroy         json.Number   `json:"spOnDestroy"`
	TargetType          json.Number   `json:"targetType"`
	TargetPattern       json.Number   `json:"targetPattern"`
	ConditionType       json.Number   `json:"conditionType"`
	HighCondition       json.Number   `json:"highCondition"`
	ConditionValue      json.Number   `json:"conditionValue"`
	MidCondition        json.Number   `json:"midCondition"`
	LowCondition        json.Number   `json:"lowCondition"`
	Bloodsucking        json.Number   `json:"bloodsucking"`
	PassiveChance       json.Number   `json:"passiveChance"`
	StartBonus          json.Number   `json:"startBonus"`
	LvBonus             json.Number   `json:"lvBonus"`
	ActionEffSound      json.Number   `json:"actionEffSound"`
	ActionEffWithFire   json.Number   `json:"actionEffWithFire"`
	ActionSoundDelay    json.Number   `json:"actionSoundDelay"`
	ActionSound         json.Number   `json:"actionSound"`
	FireSound           string        `json:"fireSound"`
	HitSoundtype        json.Number   `json:"hitSoundtype"`
	HitSound            string        `json:"hitSound"`
	HitSoundDelay       json.Number   `json:"hitSoundDelay"`
	BeHitSound          string        `json:"beHitSound"`
	BeMissSound         string        `json:"beMissSound"`
	ReturnMotion        string        `json:"returnMotion"`
	AtkDmIgDfn          json.Number   `json:"atkDmIgDfn"`
}

type ExpectedGKSkillUpgrade struct {
	Key                 json.Number `json:"key"`
	AttackDamage        json.Number `json:"attackDamage"`
	Healing             json.Number `json:"healing"`
	SpVal               json.Number `json:"spVal"`
	DotDamage           json.Number `json:"dotDamage"`
	AttackDamageBonus   string      `json:"attackDamageBonus"`
	RecvHealBonus       json.Number `json:"recvHealBonus"`
	MaxHealthBonus      json.Number `json:"maxHealthBonus"`
	SpeedBonus          json.Number `json:"speedBonus"`
	AccuracyBonus       json.Number `json:"accuracyBonus"`
	LuckBonus           json.Number `json:"luckBonus"`
	DefenseBonus        json.Number `json:"defenseBonus"`
	CriticalChanceBonus json.Number `json:"criticalChanceBonus"`
	CriticalDamageBonus string      `json:"criticalDamageBonus"`
	ShieldBonus         json.Number `json:"shieldBonus"`
	AtkDmIgDfn          json.Number `json:"atkDmIgDfn"`
	ShldCn              json.Number `json:"shldCn"`
	RmStTp              json.Number `json:"rmStTp"`
	RmStTn              json.Number `json:"rmStTn"`
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

func Test1() {
	filePath := "./json/Regulation/SkillDataTable.json"
	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.WriteErr(err)
		eMess := fmt.Sprintf("Somethings wrong with reading %s!", filePath)
		fmt.Println(eMess)
		return
	}

	var eS []ExpectedGKSkill
	if err := json.Unmarshal(res, &eS); err != nil {
		t.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(res, &ExpectedGKSkill)")
		return
	} else {
		t.WriteLog("Successful: Regulation/SkillDataTable.json Unmarshal to ExpectedGKSkill")
	}

}

func Test2() {
	filePath := "./json/Regulation/SkillUpgradeDataTable.json"
	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.WriteErr(err)
		eMess := fmt.Sprintf("Somethings wrong with reading %s!", filePath)
		fmt.Println(eMess)
		return
	}

	var eSU []ExpectedGKSkillUpgrade
	if err := json.Unmarshal(res, &eSU); err != nil {
		t.WriteErr(err)
		fmt.Println("Error with: json.Unmarshal(res, &ExpectedGKSkillUpgrade)")
		return
	} else {
		t.WriteLog("Successful: Regulation/SkillUpgradeDataTable.json Unmarshal to ExpectedGKSkillUpgrade")
	}

	// return eSU
}
