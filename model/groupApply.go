// @Title  GroupApply
// @Description  定义用户组申请
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GroupApply			定义用户组申请
type GroupApply struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                      // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                         // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`   // 用户外键
	GroupId   uuid.UUID `json:"group_id" gorm:"type:char(36);index:idx_groupId;not null"` // 用户组外键
	Condition bool      `json:"condition" gorm:"type:boolean;not null"`                   // 申请状态
	Content   string    `json:"content" gorm:"type:text;not null"`                        // 内容
	ResLong   string    `json:"res_long" gorm:"type:text"`                                // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text"`                               // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (groupApply *GroupApply) BeforeCreate(scope *gorm.DB) error {
	groupApply.ID = uuid.NewV4()
	return nil
}
