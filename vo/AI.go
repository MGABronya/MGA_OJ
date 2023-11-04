// @Title  AI
// @Description  AI模板
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
package vo

// AIRequest			接收一个AI模板
type AIRequest struct {
	Characters string `json:"characters"` // 人设
	Reply      bool   `json:"reply"`      // 是否回复自己
	Prologue   string `json:"prologue"`   // 开场白
}
