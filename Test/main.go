package main

import (
	"MGA_OJ/common"
	"MGA_OJ/controller"
	"MGA_OJ/model"
	"MGA_OJ/selfInspection"
	"MGA_OJ/util"
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
	selfInspection.TestInspection()
	InitConfig()
	common.InitDB()
	client0 := common.InitRedis(0)
	defer client0.Close()
	common.InitDocker()
	// TODO 生成公用锁
	rw := &sync.RWMutex{}
	// TODO 将锁赋予心跳控制器
	controller.HeartRW = rw
	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	// TODO 发送闲置状态心跳
	go func() {
		redis := common.GetRedisClient(0)
		ctx := context.Background()
		for {
			time.Sleep(1 * time.Second)
			// TODO 发出心跳
			rw.Lock()
			heart := model.Heart{DockerId: common.DockerId, Condition: "Waiting", TimesTamp: model.Time(time.Now())}
			v, _ := json.Marshal(heart)
			redis.Publish(ctx, "heart", v)
			rw.Unlock()
		}
	}()
	// TODO 心跳检测员
	go controller.HeartCount()
	panic(r.Run())
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
	viper.AddConfigPath(examDir + "./../config")
	err := viper.ReadInConfig()
	// TODO 如果发生错误，终止程序
	if err != nil {
		panic(err)
	}
}
