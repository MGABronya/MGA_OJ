// @Title  competitionLabel
// @Description  定义比赛的标签
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// CompetitionLabel			定义比赛标签
type CompetitionLabel struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // id
	CreatedAt     Time      `json:"created_at" gorm:"type:timestamp"`                                    // 创建日期
	UpdatedAt     Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 更新日期
	Label         string    `json:"label" gorm:"type:char(36);index:label;not null"`                     // 标签
	CompetitionId uuid.UUID `json:"compeition_id" gorm:"type:char(36);index:idx_compeition=Id;not null"` // 竞赛外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (competitionLabel *CompetitionLabel) BeforeCreate(scope *gorm.DB) error {
	competitionLabel.ID = uuid.NewV4()
	return nil
}
