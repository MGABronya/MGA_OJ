// @Title  Competition
// @Description  定义比赛
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Competition		定义比赛
type Competition struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                                    // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;"`                      // 用户外键
	GroupId   uuid.UUID `json:"group_id" gorm:"type:char(36);index:idx_groupId;"`                    // 用户组外键
	StartTime Time      `json:"start_time" gorm:"type:timestamp;not null"`                           // 起始时间
	EndTime   Time      `json:"end_time" gorm:"type:timestamp;not null"`                             // 终止时间
	Title     string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 标题
	Content   string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 内容
	ResLong   string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
	PasswdId  uuid.UUID `json:"passwd_id" gorm:"type:char(36)"`                                      // 密码id
	HackTime  Time      `json:"hack_time" gorm:"type:timestamp"`                                     // hack时间
	HackScore uint      `json:"hack_score" gorm:"type:uint"`                                         // hack分数
	HackNum   uint      `json:"hack_num" gorm:"type:uint"`                                           // hack分数封顶
	Type      string    `json:"type" gorm:"type:char(20)"`                                           // 比赛类型
	LessNum   uint      `json:"less_num" gorm:"type:uint"`                                           // 最低小组人数
	UpNum     uint      `json:"up_num" gorm:"type:uint"`                                             // 最高小组人数
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (competition *Competition) BeforeCreate(scope *gorm.DB) error {
	competition.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于比赛删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (c *Competition) BeforDelete(tx *gorm.DB) (err error) {
	// TODO 删除密码
	tx.Where("id = (?)", c.PasswdId).Delete(&Passwd{})

	tx = tx.Where("competition_id = (?)", c.ID)

	// TODO 删除比赛成员
	tx.Delete(&CompetitionMember{})

	// TODO 删除比赛排名
	tx.Delete(&CompetitionRank{})

	// TODO 删除比赛标签
	tx.Delete(&CompetitionLabel{})

	// TODO 删除题目
	tx.Delete(&ProblemNew{})

	// TODO 删除通告
	tx.Delete(&Notice{})

	return
}
