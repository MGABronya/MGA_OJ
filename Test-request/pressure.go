package test

import "time"

type User struct {
	id      string
	commits []Commit
}

type Commit struct {
	created_at time.Time
	code       string
	language   string
	problem    Problem
}

type Problem struct {
	input        []string
	time_limit   uint
	memory_limit uint
}

type Record struct {
	spand      int64
	userId     string
	created_at int64
}

var Records []Record

func (user User) Do(start time.Time, end time.Time) {
	for i := range user.commits {
		go func(i int) {
			<-time.NewTimer(user.commits[i].created_at.Sub(start)).C
			for j := range user.commits[i].problem.input {
				now := time.Now()
				JudgeRun(user.commits[i].language, user.commits[i].code, user.commits[i].problem.input[j], user.commits[i].problem.memory_limit, user.commits[i].problem.time_limit)
				spand := now.UnixMilli() - time.Now().UnixMilli()
				Records = append(Records, Record{
					spand,
					user.id,
					now.Sub(start).Milliseconds(),
				})
			}
		}(i)
	}
	<-time.NewTimer(end.Sub(start)).C
}
