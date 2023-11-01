// @Title  Heart
// @Description  定义Heart
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

// Heart		定义心跳
type Heart struct {
	DockerId  string `json:"docker_id" gorm:"type:char(36);not null"` // 容器id
	Condition string `json:"condition" gorm:"type:char(36);not null"` // 容器状态
	TimesTamp Time   `json:"times_tamp" gorm:"type:timestamp"`        // 心跳时间
}
