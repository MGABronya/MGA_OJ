// @Title  RankList
// @Description  滚榜
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

import uuid "github.com/satori/go.uuid"

// RankList			定义滚榜
type RankList struct {
	MemberId uuid.UUID `json:"member_id"`
}
