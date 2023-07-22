// @Title  badge
// @Description  定义徽章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Badge			定义徽章
type Badge struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // 徽章的id
	Name        string    `json:"name" gorm:"type:char(50);not null;unique"`              // 名称
	Description string    `json:"description" gorm:"type:text;not null;"`                 // 描述
	ResLong     string    `json:"res_long" gorm:"type:text;"`                             // 备用长文本
	ResShort    string    `json:"res_short" gorm:"type:text;"`                            // 备用短文本
	Condition   string    `json:"condition" gorm:"type:text;not null;"`                   // 获取条件
	Iron        int       `json:"iron" gorm:"type:int;not null;"`                         // 铁
	Copper      int       `json:"copper" gorm:"type:int;not null;"`                       // 铜
	Silver      int       `json:"sliver" gorm:"type:int;not null;"`                       // 银
	Gold        int       `json:"gold" gorm:"type:int;not null;"`                         // 金
	File        string    `json:"file" grom:"type:text;not null;"`                        // 文件
	UserId      uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	CreatedAt   Time      `json:"created_at" gorm:"type:timestamp"`                       // 徽章的创建日期
	UpdatedAt   Time      `json:"updated_at" gorm:"type:timestamp"`                       // 徽章的更新日期
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (badge *Badge) BeforeCreate(scope *gorm.DB) error {
	badge.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于文章删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (b *Badge) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("badge_id = (?)", b.ID)

	// TODO 删除用户勋章
	tx.Delete(&UserBadge{})

	// TODO 删除勋章行为映射
	tx.Delete(&BadgeBehavior{})

	return
}
