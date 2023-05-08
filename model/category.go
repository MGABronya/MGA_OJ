// @Title  category
// @Description  定义分类
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Category			定义分类
type Category struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`   // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`      // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`      // 更新日期
	Name      string    `json:"name" gorm:"type:varchar(64);not null"` // 分类名称
	Content   string    `json:"content" gorm:"type:text;not null"`     // 分类描述内容
	ResLong   string    `json:"res_long" gorm:"type:text"`             // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text"`            // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (category *Category) BeforeCreate(scope *gorm.DB) error {
	category.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于分类删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (c *Category) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("category_id = ?", c.ID)

	// TODO 删除相关文章
	tx.Delete(&Article{})

	return
}
