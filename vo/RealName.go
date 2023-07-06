// @Title  RealName
// @Description  定义实名
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// RealNameRequest			定义实名
type RealNameRequest struct {
	Name      string `json:"name"`       // 名字
	StudentId string `json:"student_id"` // 学号
}
