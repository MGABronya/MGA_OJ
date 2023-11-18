// @Title  File
// @Description  用于描述文件
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

import "time"

// File			定义文件
type File struct {
	Name          string    `json:"name"`          // 文件名称
	Path          string    `json:"path"`          // 文件所在路径
	Type          string    `json:"type"`          // 文件类型
	LastWriteTime time.Time `json:"lastWriteTime"` // 最后修改时间
	Size          int64     `json:"size"`          // 文件大小
}
