package guard

import (
	"MGA_OJ/common"
	"MGA_OJ/util"
	"os"

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
	InitConfig()
	// TODO 容器信息初始化
	common.InitDocker()
	// TODO redis信息初始化
	client0 := common.InitRedis(0)
	defer client0.Close()
	// TODO 检测心跳
	Guard()
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
	viper.AddConfigPath(examDir + "/config")
	err := viper.ReadInConfig()
	// TODO 如果发生错误，终止程序
	if err != nil {
		panic(err)
	}
}
