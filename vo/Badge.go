// @Title  badge
// @Description  定义徽章
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

// BadgeRequest			定义徽章
type BadgeRequest struct {
	Name        string `json:"name"`        // 名称
	Description string `json:"description"` // 描述
	ResLong     string `json:"res_long"`    // 备用长文本
	ResShort    string `json:"res_short"`   // 备用短文本
	Condition   string `json:"condition"`   // 获取条件
	Iron        int    `json:"iron"`        // 铁
	Copper      int    `json:"copper"`      // 铜
	Silver      int    `json:"sliver"`      // 银
	Gold        int    `json:"gold"`        // 金
	File        string `json:"file"`        // 文件
}
