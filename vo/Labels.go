// @Title  Labels
// @Description  接收一组标签
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

// LabelsRequest			接收一组标签
type LabelsRequest struct {
	Labels []string `form:"labels"` // 标签
}