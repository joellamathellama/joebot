package general

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"joebot/rds"
	"joebot/tools"
	"time"
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
					rr := rds.RedisLRange(rds.RC, "gkshootout")
					if len(rr) == 0 {
						return
					}
					if th == 9 || th == 21 {
						// bullet shootout
						for i := 0; i < len(rr); i++ {
							if _, err := s.ChannelMessageSend(rr[i], "Bullet Shootout!!!"); err != nil {
								tools.WriteErr(err)
								fmt.Println(err)
							}
						}
					} else if th == 3 || th == 15 {
						// gold shootout
						for i := 0; i < len(rr); i++ {
							if _, err := s.ChannelMessageSend(rr[i], "Gold Shootout!!!"); err != nil {
								tools.WriteErr(err)
								fmt.Println(err)
							}
						}
					} else {
						// hourly tick
						fmt.Println("Hourly Tick")
					}
				} else if ts == 0 {
					// minute tick
					// rr := rds.RedisLRange(rds.RC, "gkshootout")
					// if len(rr) == 0 {
					// 	return
					// }
					// for i := 0; i < len(rr); i++ {
					// 	if _, err := s.ChannelMessageSend(rr[i], "Minute Tick"); err != nil {
					// 		tools.WriteErr(err)
					// 		fmt.Println(err)
					// 	}
					// }
				}
			case <-quit:
				ticker.Stop()
				fmt.Println("Stopped the ticker!")
				return
			}
		}
	}()
}
