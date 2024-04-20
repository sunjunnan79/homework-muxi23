package main

import (
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 设置 Kafka 消费者配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 连接 Kafka
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Failed to close consumer: %v", err)
		}
	}()

	// 订阅主题
	topic := "test-topic"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}

	// 捕获 SIGINT 和 SIGTERM 信号以优雅地关闭消费者
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 接收和处理消息
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v\n", err)
		case <-signals:
			log.Println("Shutting down...")
			return
		}
	}
}
