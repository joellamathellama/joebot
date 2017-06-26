package routes

import (
	"joebot/rds"
	"joebot/tools"
)

var (
	alarmList []string
)

func CreateAlarmList() {
	alarmList = append(alarmList, "gkshootout")
}

func setAlarm(cID string, name string) (res string) {
	res = "Alarm set to this channel"
	if tools.StringInSlice(name, alarmList) {
		rds.RedisLRem(rds.RC, name, cID)
		rds.RedisLPush(rds.RC, name, cID)
	} else {
		res = "Alarm does not exist"
	}
	return
}

func removeAlarm(cID string, name string) (res string) {
	res = "Alarm removed from this channel"
	if tools.StringInSlice(name, alarmList) {
		rds.RedisLRem(rds.RC, name, cID)
	} else {
		res = "Alarm does not exist"
	}
	return
}
