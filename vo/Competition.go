// @Title  CompetitionRequest
// @Description  定义比赛
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package vo

import (
	"MGA_OJ/model"

	uuid "github.com/satori/go.uuid"
)

// CompetitionRequest		定义比赛
type CompetitionRequest struct {
	SetId     uuid.UUID  `json:"set_id"`     // 表单外键
	StartTime model.Time `json:"start_time"` // 起始时间
	EndTime   model.Time `json:"end_time"`   // 终止时间
	Type      string     `json:"type"`       // 比赛类型
}

// CompdeitionUpdate		定义比赛更新
type CompetitionUpdate struct {
	StartTime model.Time `json:"start_time"` // 起始时间
	EndTime   model.Time `json:"end_time"`   // 终止时间
}
