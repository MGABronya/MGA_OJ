// @Title  notice
// @Description  定义通告
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Notice			定义通告
type Notice struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                  // 文章的id
	CreatedAt     Time      `json:"created_at" gorm:"type:timestamp"`                                     // 文章的创建日期
	UpdatedAt     Time      `json:"updated_at" gorm:"type:timestamp"`                                     // 文章的更新日期
	UserId        uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`               // 用户外键
	CompetitionId uuid.UUID `json:"competition_id" gorm:"type:char(36);index:idx_competitionId;not null"` // 分类外键
	Title         string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"`  // 文章的标题
	Content       string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`    // 文章的内容
	ResLong       string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`            // 备用长文本
	ResShort      string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (notice *Notice) BeforeCreate(scope *gorm.DB) error {
	notice.ID = uuid.NewV4()
	return nil
}
