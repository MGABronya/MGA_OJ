// @Title  case
// @Description  定义提交记录中的样例
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

// Case			定义提交记录中的样例
type Case struct {
	RecordId uint `json:"record_id" gorm:"type:uint;index:idx_recordId;not null"` // 提交记录外键
	ID       uint `json:"id" gorm:"type:uint;not null"`                           // 表示第几个测试
	Time     uint `json:"time" gorm:"type:uint;not null"`                         // 表示测试使用时间
	Memory   uint `json:"memory" gorm:"type:uint;not null"`                       // 表示测试使用空间
}
