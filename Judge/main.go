package main

import (
	"MGA_OJ/common"
	"MGA_OJ/consumer"
	"MGA_OJ/selfInspection"
	"MGA_OJ/util"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

// @title    main
// @description   程序入口，完成一些初始化工作后将开始监听
// @auth      MGAronya             2022-9-16 10:49
// @param     void			没有入参
// @return    void			没有回参
func main() {
	// TODO 打印MGAronya字符串
	util.MgaronyaPrint()
	// TODO 自检程序启动
	selfInspection.JudgeInspection()
	InitConfig()
	common.InitDB()
	client0 := common.InitRedis(0)
	defer client0.Close()
	rw := &sync.RWMutex{}
	consumer.NewGroup(rw)
	consumer.NewSingle(rw)
	consumer.NewJudge(rw)
	common.InitRabbitmq()
	RabbitMQ := common.GetRabbitMq()
	log.Println("Consumer examing...")
	RabbitMQ.ConsumeSimple()
}

// @title    InitConfig
// @description   读取配置文件并完成初始化
// @auth      MGAronya             2022-9-16 10:49
// @param     void			没有入参
// @return    void			没有回参
func InitConfig() {
	examDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(examDir + "/../config")
	err := viper.ReadInConfig()
	// TODO 如果发生错误，终止程序
	if err != nil {
		panic(err)
	}
}
