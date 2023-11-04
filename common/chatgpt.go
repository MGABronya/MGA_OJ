// @Title  chatgpt
// @Description  该文件用于初始化bing搜索api
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package common

import (
	"github.com/spf13/viper"
)

// ChatGPTUrl			chat-gpt3.5 API的端点URL
var ChatGPTUrl string

// ChatGPTKey			包含API密钥的头部信息
var ChatGPTKey string

// @title    InitChatGPT
// @description   从配置文件中读取ChatGPT的api相关参数
// @auth      MGAronya             2022-9-16 10:07
// @param     void        void         没有入参
// @return    void        void         没有返回值
func InitChatGPT() {
	ChatGPTUrl = viper.GetString("chatgpt.url")
	ChatGPTKey = viper.GetString("chatgpt.apikey")
}
