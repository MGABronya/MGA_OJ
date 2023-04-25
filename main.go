package main

import (
	"MGA_OJ/common"
	"MGA_OJ/controller"
	"MGA_OJ/selfInspection"
	"MGA_OJ/timer"
	"MGA_OJ/util"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title    main
// @description   程序入口，完成一些初始化工作后将开始监听
// @auth      MGAronya（张健）             2022-9-16 10:49
// @param     void			没有入参
// @return    void			没有回参
func main() {
	// TODO 打印MGAronya字符串
	util.MgaronyaPrint()
	// TODO 自检程序启动
	selfInspection.MainInspection()
	InitConfig()
	common.InitDB()
	client0 := common.InitRedis(0)
	defer client0.Close()
	common.InitRabbitmq()
	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	// TODO 定时任务
	go timer.TimedTask("The hot data is being updated...", timer.HotStatics, 0, 0, 0)
	go timer.TimedTask("Form ranking calculation...", timer.SetRank, 4, 0, 0)
	// TODO 及时单人匹配比赛
	go controller.CompetitionRandomSingleGo()
	// TODO 及时小组匹配比赛
	go controller.CompetitionRandomGroupGo()

	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

// @title    InitConfig
// @description   读取配置文件并完成初始化
// @auth      MGAronya（张健）             2022-9-16 10:49
// @param     void			没有入参
// @return    void			没有回参
func InitConfig() {
	examDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(examDir + "/config")
	err := viper.ReadInConfig()
	// TODO 如果发生错误，终止程序
	if err != nil {
		panic(err)
	}
}
