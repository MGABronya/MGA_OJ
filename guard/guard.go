package main

import (
	"MGA_OJ/Docker"
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"context"
	"encoding/json"
	"log"
	"time"
)

// DockerMap    定义了docker映射
var DockerMap map[string]Interface.GuardInterface = map[string]Interface.GuardInterface{
	"test":  Docker.NewTest(),
	"judge": Docker.NewJudge(),
}

// @title    Guard
// @description   持续检测目标心跳
// @auth      MGAronya             2022-9-16 10:49
// @param     void			没有入参
// @return    void			没有回参
func Guard() {
	redis := common.GetRedisClient(0)
	ctx := context.Background()
	docker := DockerMap[common.Type]
	// TODO 订阅消息
	pubSub := redis.Subscribe(ctx, "heart")
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()
	// TODO 计时器
	t := time.NewTimer(time.Second * 30)
	// TODO 如果300s内没有响应，重启容器
	go func() {
		for {
			// TODO 时间到了
			<-t.C
			log.Println("容器失联")
			// TODO 停止容器
			cmd := docker.Stop(common.DockerId)
			if err := cmd.Run(); err != nil {
				log.Println("停止容器失败:" + err.Error())
			}

			// TODO 移除容器
			cmd = docker.RM(common.DockerId)
			if err := cmd.Run(); err != nil {
				log.Println("移除容器失败:" + err.Error())
			}

			// TODO 创建容器
			cmd = docker.Run(common.DockerId, common.CPU, common.PostMap)
			if err := cmd.Run(); err != nil {
				log.Println("创建容器失败:" + err.Error())
			}
			log.Println("容器已重建")
			// TODO 更新计时
			t.Reset(time.Second * 300)
		}
	}()
	// TODO 监听消息
	for msg := range ch {
		var heart model.Heart
		json.Unmarshal([]byte(msg.Payload), &heart)
		// TODO 查看id是否对上
		if heart.DockerId != common.DockerId {
			continue
		}
		// TODO 更新计时
		t.Reset(time.Second * 300)
	}
}
