// @Title  translation
// @Description  该文件用于初始化bing翻译api
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package common

import (
	"github.com/spf13/viper"
)

var TraslationSubscriptionKey []string

// @title    InitTranslation
// @description   从配置文件中读取bing翻译api相关参数
// @auth      MGAronya             2022-9-16 10:07
// @param     void        void         没有入参
// @return    void        void         没有返回值
func InitTranslation() {
	TraslationSubscriptionKey = append(TraslationSubscriptionKey, viper.GetString("translation.subscriptionKey1"))
	TraslationSubscriptionKey = append(TraslationSubscriptionKey, viper.GetString("translation.subscriptionKey2"))
}
