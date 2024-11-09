package main

import (
	"Insomnia/app/api/routers"
	"Insomnia/app/infrastructure/kafka"
	"Insomnia/app/infrastructure/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title 不眠夜
// @version 1.0
// @description 一个匿名熬夜论坛
func main() {
	//启动gin的engine
	engine := gin.Default()
	routers.Load(engine)
	//加载中间件
	middlewares.Load(engine)
	//加载消费者组
	go ConsumerGroup()
	if err := engine.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

func ConsumerGroup() {
	// 创建 Kafka 实例
	topics := []string{"cache"}
	group := "cache-group"
	key := "cache-key"
	kafkaClient := kafka.NewKafka(topics, group, key)

	// 创建并启动消费者组
	closeConsumer := kafkaClient.CreateConsumerToGroup()
	defer closeConsumer()

	// 监听信号以优雅地退出程序
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
}
