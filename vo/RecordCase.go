// @Title  RecorkList
// @Description  提交频道
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

// RecordCase			定义提交频道消息
type RecordCase struct {
	Condition string `json:"condition"`
	CaseId    uint   `json:"case_id"`
}
