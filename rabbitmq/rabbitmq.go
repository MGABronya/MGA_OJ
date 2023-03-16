package rabbitMq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// 连接信息amqp://MGAronya:MGAronya@127.0.0.1:5672/MGAronya这个信息是固定不变的amqp://事固定参数后面两个是用户名密码ip地址端口号Virtual Host
const MQURL = "amqp://MGAronya:MGAronya@127.0.0.1:5672/MGAronya"

var max_run int = 2

var ch chan struct{} = make(chan struct{}, max_run)

// RabbitMQ		定义了rabbitMQ结构体
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机名称
	Key       string // bind Key 名称
	Mqurl     string // 连接信息
}

// @title    NewRabbitMQ
// @description  创建结构体实例
// @auth      MGAronya（张健）       2022-11-25 12:20
// @param    queueName string, exchange string, key string	接收队列名称，交换机名称，bind Key名称和连接信息
// @return   *RabbitMQ			返回一个rabbitMQ实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// @title    Destory
// @description  断开channel 和 connection
// @auth      MGAronya（张健）       2022-11-25 12:20
// @param    void			没有参数
// @return   void			没有返回值
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// @title    failOnErr
// @description  错误处理函数
// @auth      MGAronya（张健）       2022-11-25 12:20
// @param    err error, message string			错误信息
// @return   void								没有返回值
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// @title    NewRabbitMQSimple
// @description  创建简单模式下RabbitMQ实例
// @auth      MGAronya（张健）       2022-11-25 12:20
// @param    queueName string			队列名
// @return   *RabbitMQ					一个rabbitmq实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// TODO 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	// TODO 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+"itmq!")
	// TODO 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// @title    PublishSimple
// @description  直接模式队列生产
// @auth      MGAronya（张健）       2022-11-25 12:20
// @param    message string			信息
// @return   void					没有返回值
func (r *RabbitMQ) PublishSimple(message string) error {
	// TODO 申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞处理
		nil,   // 额外的属性
	)
	if err != nil {
		return err
	}
	// TODO 调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false, // 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

// @title    ConsumeSimple
// @description  simple 模式下消费者
// @auth      MGAronya（张健）       2022-11-25 12:20
// @param    void					没有参数
// @return   void					没有返回值
func (r *RabbitMQ) ConsumeSimple() {
	// TODO 申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞处理
		nil,   // 额外的属性
	)
	if err != nil {
		fmt.Println(err)
	}

	// TODO 接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue,用来区分多个消费者
		"",     // consumer
		true,   // auto-ack,是否自动应答
		false,  // exclusive,是否独有
		false,  // no-local,设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中的消费者
		false,  // no-wait,列是否阻塞
		nil,    // args
	)

	if err != nil {
		fmt.Println(err)
	}

	for d := range msgs {
		// TODO 在管道内放入正在运行时，道满时这里会阻塞
		ch <- struct{}{}

		// TODO 启用协程处理消息
		go func(body []byte) {
			Judge(body)
			// TODO 完成处理后，从管道中拿出一份处理消息
			<-ch
		}(d.Body)

	}
}
