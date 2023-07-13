// @Title  tag
// @Description  用于自动生成的标签
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

// Tag			定义自动生成标签
type Tag struct {
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp"`         // 创建日期
	UpdatedAt Time   `json:"updated_at" gorm:"type:timestamp"`         // 更新日期
	Tag       string `json:"tag" gorm:"type:char(36);not null;unique"` // 标签
}
