package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

func main() {
	// 设置 Kafka 生产者配置
	config := sarama.NewConfig()
	// 设置消息发送配置
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认消息
	config.Producer.Retry.Max = 3                    // 消息发送失败时的最大重试次数
	config.Producer.Return.Successes = true          // 返回成功发送的消息

	// 设置 TLS 配置
	//config.Net.TLS.Enable = true // 启用 TLS 加密通信

	// 设置 SASL 认证配置
	//config.Net.SASL.Enable = true     // 启用 SASL 认证
	//config.Net.SASL.User = "user"     // SASL 用户名
	//config.Net.SASL.Password = "pass" // SASL 密码

	// 连接 Kafka
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Failed to close producer: %v", err)
		}
	}()

	// 发送消息
	msg := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Value:     sarama.StringEncoder("Hello, Kafka!"),
		Partition: 1,
		Timestamp: time.Now().Add(60 * time.Second),
		Key:       sarama.StringEncoder("贱!"), //好孩子不要用中文做key哦
		Headers: []sarama.RecordHeader{
			{Key: []byte("header1"), Value: []byte("value1")},
			{Key: []byte("header2"), Value: []byte("value2")},
		},
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
