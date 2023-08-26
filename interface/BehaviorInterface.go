// @Title  BehaviorInterface
// @Description  该文件用于封装行为方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import (
	uuid "github.com/satori/go.uuid"
)

// BehaviorInterface			定义了行为方法
type BehaviorInterface interface {
	UserBehavior(uuid.UUID) (float64, error)  // 查看某个用户的行为统计
	PublishBehavior(float64, uuid.UUID) error // 更新行为统计，并按情况通报
	Description() string                      // 返回行为描述
}
