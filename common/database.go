// @Title  database
// @Description  该文件用于初始化mysql数据库，以及包装一个向外提供数据库的功能
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package common

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// @title    InitDB
// @description   从配置文件中读取数据库相关信息后，完成数据库初始化
// @auth      MGAronya（张健）             2022-9-16 10:07
// @param     void        void         没有入参
// @return    db        *gorm.DB         将返回一个初始化后的数据库指针
func InitDB() *gorm.DB {
	//driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	// TODO  如果未能连接到数据库，终止程序并返回错误信息
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	DB = db
	return db
}

// @title    GetDB
// @description   返回数据库的指针
// @auth      MGAronya（张健）             2022-9-16 10:08
// @param     void        void         没有入参
// @return    db        *gorm.DB         将返回一个初始化后的数据库指针
func GetDB() *gorm.DB {
	return DB
}
