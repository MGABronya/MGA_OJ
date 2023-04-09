package main

import (
	TQ "MGA_OJ/Test-request"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// 本地测试使用
func main() {
	var users []TQ.User

	var start, end time.Time

	dstFile, err := os.Create("pressure.dat")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	for i := range users {
		go users[i].Do(start, end)
	}

	<-time.NewTimer(end.Sub(start)).C

	bs, _ := json.Marshal(TQ.Records)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")

	dstFile.WriteString(out.String())

}
