// @Title  student
// @Description  定义学生
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

// Student			定义学生
type Student struct {
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp"`              // 创建日期
	UpdatedAt Time   `json:"updated_at" gorm:"type:timestamp"`              // 更新日期
	Name      string `json:"content" gorm:"type:char(36);not null"`         // 姓名
	StudentId string `json:"res_long" gorm:"type:char(36);not null;unique"` // 学号
	Major     string `json:"res_short" gorm:"type:char(36);not null"`       // 专业
	Grade     string `json:"grade" gorm:"type:char(36);not null"`           // 年级
	College   string `json:"college" gorm:"type:char(12);not null"`         // 学院
}
