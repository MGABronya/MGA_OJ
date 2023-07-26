// @Title  search
// @Description  该文件用于初始化bing搜索api
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package common

import (
	"github.com/spf13/viper"
)

var SubscriptionKey []string

// @title    InitSearch
// @description   从配置文件中读取bing搜索api相关参数
// @auth      MGAronya（张健）             2022-9-16 10:07
// @param     void        void         没有入参
// @return    void        void         没有返回值
func InitSearch() {
	SubscriptionKey = append(SubscriptionKey, viper.GetString("search.subscriptionKey1"))
	SubscriptionKey = append(SubscriptionKey, viper.GetString("search.subscriptionKey2"))
}
