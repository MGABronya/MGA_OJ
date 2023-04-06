// @Title  ConsumerInterface
// @Description  该文件用于封装消费者方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

// ConsumerInterface			定义了消费者方法
type ConsumerInterface interface {
	Handel(body []byte) // 对消息进行消费
}
