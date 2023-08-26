// @Title  ConsumerInterface
// @Description  该文件用于封装消费者方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

// ConsumerInterface			定义了消费者方法
type ConsumerInterface interface {
	Handel(body string) // 对消息进行消费
}

var ComsumerMap map[string]ConsumerInterface = map[string]ConsumerInterface{}
