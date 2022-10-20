// @Title  user
// @Description  定义用户
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import "gorm.io/gorm"

// user			定义用户
type User struct {
	gorm.Model        // gorm的模板
	Name       string `json:"name" gorm:"type:varchar(20);not null;unique"`  // 用户名称
	Email      string `json:"email" gorm:"type:varchar(50);not null;unique"` // 邮箱
	Password   string `json:"password" gorm:"size:255;not null"`             // 密码
	Icon       string `json:"icon" gorm:"type:varchar(50)"`                  // 这里的Icon存储的是图像文件的url后缀
	Blog       string `json:"blog" gorm:"type:varchar(25)"`                  // 博客
	Sex        bool   `json:"sex" gorm:"type:boolean"`                       // 性别
	Address    string `json:"address" gorm:"type:varchar(20)"`               // 地址
	Level      int    `json:"level" gorm:"type:int;not null"`                // 用户管理等级
}
