// @Title  CompetitionRequest
// @Description  定义比赛
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package vo

import (
	"MGA_OJ/model"

	uuid "github.com/satori/go.uuid"
)

// CompetitionRequest		定义比赛
type CompetitionRequest struct {
	StartTime model.Time `json:"start_time"` // 起始时间
	EndTime   model.Time `json:"end_time"`   // 终止时间
	Title     string     `json:"title"`      // 标题
	Content   string     `json:"content"`    // 内容
	ResLong   string     `json:"res_long"`   // 备用长文本
	ResShort  string     `json:"res_short"`  // 备用短文本
	HackTime  model.Time `json:"hack_time"`  // hack时间
	HackScore uint       `json:"hack_score"` // hack分数
	HackNum   uint       `json:"hack_num"`   // hack分数封顶
	GroupId   uuid.UUID  `json:"group_id"`   // 用户组
	RealName  bool       `json:"real_name"`  // 是否需要实名
	Type      string     `json:"type"`       // 类型
	LessNum   uint       `json:"less_num"`   // 最低小组人数
	UpNum     uint       `json:"up_num"`     // 最高小组人数
}
