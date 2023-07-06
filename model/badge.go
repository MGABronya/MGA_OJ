// @Title  badge
// @Description  定义徽章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// Badge			定义徽章
type Badge struct {
	Name        string    `json:"name" gorm:"type:char(50);not null;"`                    // 名称
	Description string    `json:"description" gorm:"type:text;not null;"`                 // 描述
	ResLong     string    `json:"res_long" gorm:"type:text;"`                             // 备用长文本
	ResShort    string    `json:"res_short" gorm:"type:text;"`                            // 备用短文本
	Condition   string    `json:"condition" gorm:"type:text;not null;"`                   // 获取条件
	Iron        int       `json:"iron" gorm:"type:int;not null;"`                         // 铁
	Copper      int       `json:"copper" gorm:"type:int;not null;"`                       // 铜
	Silver      int       `json:"sliver" gorm:"type:int;not null;"`                       // 银
	Gold        int       `json:"gold" gorm:"type:int;not null;"`                         // 金
	UserId      uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	CreatedAt   Time      `json:"created_at" gorm:"type:timestamp"`                       // 徽章的创建日期
	UpdatedAt   Time      `json:"updated_at" gorm:"type:timestamp"`                       // 徽章的更新日期
}