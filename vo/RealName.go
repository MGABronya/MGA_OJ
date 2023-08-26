// @Title  RealName
// @Description  定义实名
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

// RealNameRequest			定义实名
type RealNameRequest struct {
	Name      string `json:"name"`       // 名字
	StudentId string `json:"student_id"` // 学号
}
