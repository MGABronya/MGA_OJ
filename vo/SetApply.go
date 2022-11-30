// @Title  SetApply
// @Description  表单申请
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// SetApplyRequest			定义表单申请
type SetApplyRequest struct {
	Content  string `json:"content"`   // 内容
	Reslong  string `json:"res_long"`  // 备用长文本
	Resshort string `json:"res_short"` // 备用短文本
	GroupId  uint   `json:"group_id"`  // 用户组id
}
