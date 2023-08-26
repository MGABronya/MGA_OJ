// @Title  search
// @Description  该文件用于初始化bing搜索api
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package common

import (
	"github.com/spf13/viper"
)

var SearchSubscriptionKey []string

// @title    InitSearch
// @description   从配置文件中读取bing搜索api相关参数
// @auth      MGAronya             2022-9-16 10:07
// @param     void        void         没有入参
// @return    void        void         没有返回值
func InitSearch() {
	SearchSubscriptionKey = append(SearchSubscriptionKey, viper.GetString("search.subscriptionKey1"))
	SearchSubscriptionKey = append(SearchSubscriptionKey, viper.GetString("search.subscriptionKey2"))
}
