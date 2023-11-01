// @Title  docker
// @Description  该文件用于初始化bing翻译api
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package common

import (
	"github.com/spf13/viper"
)

var DockerId string

// @title    InitDocker
// @description   从配置文件中读取docker相关参数
// @auth      MGAronya             2022-9-16 10:07
// @param     void        void         没有入参
// @return    void        void         没有返回值
func InitDocker() {
	DockerId = viper.GetString("docker.id")
}
