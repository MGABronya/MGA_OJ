// @Title  Set
// @Description  定义表单
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// SetRequest			定义表单
type SetRequest struct {
	Title      string `json:"title"`       // 题目
	Content    string `json:"content"`     // 内容
	Reslong    string `json:"res_long"`    // 备用长文本
	Resshort   string `json:"res_short"`   // 备用短文本
	Topics     []int  `json:"topics"`      // 所包含的主题
	Groups     []int  `json:"groups"`      // 所包含的用户组
	AutoUpdate bool   `json:"auto_update"` // 是否每日自动更新排行
	AutoPass   bool   `json:"auto_pass"`   // 是否自动通过组加入申请
	PassNum    uint   `json:"pass_num"`    // 每组的人数限制
	PassRe     bool   `json:"pass_re"`     // 每组人员是否可以重复
}
