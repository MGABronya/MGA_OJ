package timer

import (
	"log"
	"time"
)

func TimedTask(des string, fun func(), hour int, min int, sec int) {
	for {
		now := time.Now()
		// TODO 计算下一个4:30
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, sec, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		// TODO 执行功能
		log.Println(des)
		fun()
	}
}
