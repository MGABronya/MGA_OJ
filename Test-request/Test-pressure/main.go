package main

import (
	TQ "MGA_OJ/Test-request"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var competitions []TQ.Competition

func initCommits(pathname string, competition *TQ.Competition) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return
	}
	usersMap := make(map[string]bool)
	var commits []TQ.Commit
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		fullName := pathname + "/" + fi.Name()
		file, _ := os.Open(fullName)
		// 初始化一个 csv reader，并通过这个 reader 从 csv 文件读取数据
		reader := csv.NewReader(file)
		// 设置返回记录中每行数据期望的字段数，-1 表示返回所有字段
		reader.FieldsPerRecord = -1
		// 通过 readAll 方法返回 csv 文件中的所有内容
		record, err := reader.ReadAll()
		if err != nil {
			println(fullName, err)
			return
		}

		if len(record) < 2 {
			continue
		}
		// TODO 遍历从 csv 文件中读取的所有内容
		for i := 1; i < len(record); i++ {
			var commit TQ.Commit
			commit.ProblemNum = record[i][1]
			if record[i][2][0] == 'C' {
				record[i][2] = "C"
			} else if record[i][2][0] == 'J' {
				record[i][2] = "Java"
			} else {
				record[i][2] = "Python"
			}
			commit.Language = record[i][2]
			commit.Created_at, _ = time.Parse("2006-01-02 15:04:05", record[i][4])
			commit.UserId = pathname + record[i][5]
			commit.Code = record[i][7]
			commits = append(commits, commit)
			usersMap[record[i][5]] = true
		}
	}
	fmt.Printf("Competition %s has %d participants and %d commits\n", pathname, len(usersMap), len(commits))
	(*competition).Commits = commits
}

// 本地测试使用
func main() {
	TQ.Records = make([]TQ.Record, 0)
	competitions = make([]TQ.Competition, 0)
	println("Init Competitions...")
	initCompetitions1()
	initCompetitions2()
	initCompetitions3()
	initCompetitions4()
	initCompetitions5()
	initCompetitions6()
	initCompetitions7()
	initCompetitions8()
	initCompetitions9()
	initCompetitions10()
	println("Init Commites...")
	initCommits("0", &competitions[0])
	initCommits("1", &competitions[1])
	initCommits("2", &competitions[2])
	initCommits("3", &competitions[3])
	initCommits("4", &competitions[4])
	initCommits("5", &competitions[5])
	initCommits("6", &competitions[6])
	initCommits("7", &competitions[7])
	initCommits("8", &competitions[8])
	initCommits("9", &competitions[9])

	dstFile, err := os.Create("pressure.dat")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	println("Runing...")
	for i := range competitions {
		go competitions[i].Do(5)
	}

	<-time.NewTimer(125 * time.Minute).C

	bs, _ := json.Marshal(TQ.Records)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")

	dstFile.WriteString(out.String())
	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
