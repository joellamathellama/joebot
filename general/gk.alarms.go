package general

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/rds"
	"joebot/tools"
	"time"
)

var (
	alarm_msg string
	rr        []string
)

func AlarmGKShootout(s *dg.Session) {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				dt_now := time.Now().UTC()
				th, tm, ts := time.Time.Clock(dt_now)
				if tm == 0 && ts == 0 {
					if th == 9 || th == 21 {
						// bullet shootout
						alarm_msg = fmt.Sprintf("%s%s%s%s\n",
							alarm_msg,
							":gun::gun::gun:",
							"Bullet Shootout!!!",
							":gun::gun::gun:")
					} else if th == 3 || th == 15 {
						// gold shootout
						alarm_msg = fmt.Sprintf("%s%s%s%s\n",
							alarm_msg,
							":moneybag::moneybag::moneybag:",
							"Gold Shootout!!!",
							":moneybag::moneybag::moneybag:")
					}
					// hourly tick
					// alarm_msg = fmt.Sprintf("%s%s\n",
					// 	alarm_msg,
					// 	"Hourly Tick")
					// fmt.Println("Hourly Tick")
				}
				if ts == 0 {
					// minute tick
					// rr = rds.RedisLRange(rds.RC, "gkshootout")
					// if len(rr) == 0 {
					// 	return
					// }
					// alarm_msg = fmt.Sprintf("%s%s%s\n",
					// 	alarm_msg,
					// 	"Minute Tick",
					// 	":gun::moneybag: test emoji :gun::moneybag:")
					// fmt.Println("Minute Tick")
				}
				// seconds tick
				if alarm_msg != "" {
					rr = rds.RedisLRange(rds.RC, "gkshootout")
					if len(rr) == 0 {
						return
					}
					for i := 0; i < len(rr); i++ {
						if _, err := s.ChannelMessageSend(rr[i], alarm_msg); err != nil {
							tools.WriteErr(err)
							fmt.Println(err)
						}
					}
				}
				alarm_msg = ""
			case <-quit:
				ticker.Stop()
				fmt.Println("Stopped the ticker!")
				return
			}
		}
	}()
}
