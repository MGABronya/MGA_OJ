// @Title  rabbitmq
// @Description  该文件用于初始化rabbitmq连接，以及包装一个向外提供rabbitmq连接的功能
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package common

import (
	"MGA_OJ/Interface"
	"MGA_OJ/model"
	"context"
	"sync"
	"time"

	"MGA_OJ/vo"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// RabbitMQ		定义了rabbitMQ结构体
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机名称
	Key       string // bind Key 名称
	Mqurl     string // 连接信息
}

var rabbitmq *RabbitMQ

// @title    Destory
// @description  断开channel 和 connection
// @auth      MGAronya       2022-11-25 12:20
// @param    void			没有参数
// @return   void			没有返回值
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// @title    failOnErr
// @description  错误处理函数
// @auth      MGAronya       2022-11-25 12:20
// @param    err error, message string			错误信息
// @return   void								没有返回值
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// @title    PublishSimple
// @description  直接模式队列生产
// @auth      MGAronya       2022-11-25 12:20
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
// @auth      MGAronya       2022-11-25 12:20
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

	// TODO 此处为心跳声明锁
	var rw *sync.RWMutex = &sync.RWMutex{}
	redis := GetRedisClient(0)
	ctx := context.Background()

	// TODO 发送闲置状态心跳
	go func() {
		for {
			time.Sleep(1 * time.Second)
			// TODO 发出心跳
			rw.Lock()
			heart := model.Heart{DockerId: DockerId, Condition: "Waiting", TimesTamp: model.Time(time.Now())}
			v, _ := json.Marshal(heart)
			redis.Publish(ctx, "heart", v)
			rw.Unlock()
		}
	}()

	for d := range msgs {
		log.Println("consumer:", d.Body)
		// TODO 进行消费
		var recordRabbit vo.RecordRabbit
		if err := json.Unmarshal(d.Body, &recordRabbit); err != nil {
			log.Println("Error Json Unmarshal:", err)
			continue
		}
		if _, ok := Interface.ComsumerMap[recordRabbit.Type]; !ok {
			log.Println("Error Record Type:", recordRabbit.Type)
			continue
		}
		// TODO 心跳检测，开始忙碌
		rw.Lock()
		heart := model.Heart{DockerId: DockerId, Condition: "Running", TimesTamp: model.Time(time.Now())}
		v, _ := json.Marshal(heart)
		redis.Publish(ctx, "heart", v)
		Interface.ComsumerMap[recordRabbit.Type].Handel(recordRabbit.RecordId.String())
		// TODO 心跳检测，忙碌完成
		heart = model.Heart{DockerId: DockerId, Condition: "Finish", TimesTamp: model.Time(time.Now())}
		v, _ = json.Marshal(heart)
		redis.Publish(ctx, "heart", v)
		rw.Unlock()
	}
}

// @title    InitRabbitmq
// @description   从配置文件中读取rabbitmq相关信息后，完成rabbitmq初始化
// @auth      MGAronya             2022-9-16 10:07
// @param     void        void         没有入参
// @return    db        *Rabbit         将返回一个初始化后的rabbitmq指针
func InitRabbitmq() *RabbitMQ {
	host := viper.GetString("rabbitmq.host")
	port := viper.GetString("rabbitmq.port")
	username := viper.GetString("rabbitmq.username")
	password := viper.GetString("rabbitmq.password")
	virtual := viper.GetString("rabbitmq.virtual")
	args := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		username,
		password,
		host,
		port,
		virtual,
	)
	rabbitmq = &RabbitMQ{QueueName: "MGAronya", Exchange: "", Key: "", Mqurl: args}
	var err error
	// TODO 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// TODO 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// @title    GetRabbitMq
// @description   返回rabbitmq的指针
// @auth      MGAronya             2022-9-16 10:08
// @param     void        void         没有入参
// @return    rabbitmq        *RabbitMQ         将返回一个初始化后的RabbitMQ指针
func GetRabbitMq() *RabbitMQ {

	return rabbitmq
}
