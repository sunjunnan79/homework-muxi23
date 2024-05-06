package kafka

import (
	. "Insomnia/app/infrastructure/helper"
	"Insomnia/app/infrastructure/redis"
	"Insomnia/app/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Kafka struct {
	brokers           []string
	topics            []string
	startOffset       int64
	version           string
	ready             chan bool
	group             string
	channelBufferSize int
	assignor          string
	key               string
}

var brokers = []string{"localhost:9092"}
var assignor = "range"

func NewKafka(topics []string, group string, key string) *Kafka {
	return &Kafka{
		brokers:           brokers,
		topics:            topics,
		group:             group,
		channelBufferSize: 1000,
		ready:             make(chan bool),
		version:           "2.8.0",
		assignor:          assignor,
		key:               key,
	}
}

// CreateCacheProducer key表示要发送的对象(缓存的对象),num表示要缓存的指定帖子页数
func (k *Kafka) CreateCacheProducer(uid string) {
	// 设置 Kafka 生产者的配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	// 创建 Kafka 生产者
	producer, err := sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	//创建消息
	message := &sarama.ProducerMessage{
		Topic: "cache",
		Key:   sarama.StringEncoder(k.key),
		Value: sarama.StringEncoder(uid),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}

func (k *Kafka) CreateConsumerToGroup() func() {
	log.Infoln("kafka init...")
	//解析卡夫卡的版本
	version, err := sarama.ParseKafkaVersion(k.version)
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}
	//初始化配置
	config := sarama.NewConfig()
	config.Version = version

	// 分区分配策略
	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	//读取措施为从最早的开始
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	//设置消费者组的channel容量,即一次最多读取几个
	config.ChannelBufferSize = k.channelBufferSize // channel长度

	// 创建client
	newClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	// 获取所有的topic,似乎没什么意义
	//topics, err := newClient.Topics()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Info("topics: ", topics)

	// 根据client和kafka指定的组别创建consumerGroup
	client, err := sarama.NewConsumerGroupFromClient(k.group, newClient)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}

	//用于控制goroutine,不是很懂但是用的蛮多的
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//死循环获取消息
		for {
			//消费者setup失败
			if err := client.Consume(ctx, k.topics, k); err != nil {
				// 当setup失败的时候，error会返回到这里
				log.Errorf("Error from consumer: %v", err)
				return
			}
			//检查ctx是否出错,一点也不懂
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			//把卡夫卡的ready状态清除
			k.ready = make(chan bool)
		}
	}()
	//如果k.ready接受到消息的话就执行
	<-k.ready
	log.Infoln("Sarama consumer up and running!...")
	// 保证在系统退出时，通道里面的消息被消费
	return func() {
		log.Info("kafka close")
		cancel()
		wg.Wait()
		if err = client.Close(); err != nil {
			log.Errorf("Error closing client: %v", err)
		}
	}
}

func (k *Kafka) Setup(sarama.ConsumerGroupSession) error {
	log.Info("setup")
	return nil
}

func (k *Kafka) Cleanup(sarama.ConsumerGroupSession) error {
	log.Info("cleanup")
	return nil
}

func (k *Kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batchSize := 10 // 定义每次处理的消息数量
	//定义一个用于临时存储消息数据的
	messages := make([]*sarama.ConsumerMessage, 0, batchSize)

	for {
		select {
		case <-session.Context().Done():
			return nil
		default:
			// 尝试从claim.Messages()中获取消息，直到达到批量处理阈值或者没有更多消息
			for i := 0; i < batchSize; i++ {
				select {
				case message, ok := <-claim.Messages():
					if !ok {
						break
					}
					messages = append(messages, message)
				case <-session.Context().Done():
					return nil
				}
			}

			// 处理批量消息
			for _, message := range messages {
				exist, err := redis.ExistResp(string(message.Value))
				if exist != 1 && err == nil {
					if string(message.Key) == "posts" {
						p := &service.PostService{}
						posts, err := p.GetPosts(string(message.Value))
						if err != nil {
							Danger(err, "加载缓存时取时获取post数据失败")
						}
						jsonString, err := json.Marshal(posts)
						if err != nil {
							Danger(err, "格式转换失败")
						}
						err = redis.SetResp(string(message.Key)+string(message.Value), jsonString)
						if err != nil {
							Danger(err, "加载缓存时存储数据失败")
						}
						fmt.Println("posts缓存存储成功")
					} else if string(message.Key) == "reposts" {
						p := &service.RePostService{}
						reposts, err := p.GetRePosts(string(message.Value))
						if err != nil {
							Danger(err, "加载缓存时取时获取repost数据失败")
						}
						jsonString, err := json.Marshal(reposts)
						if err != nil {
							Danger(err, "格式转换失败")
						}
						err = redis.SetResp(string(message.Key)+string(message.Value), jsonString)
						if err != nil {
							Danger(err, "加载缓存时取时存储数据失败")
						}
						fmt.Println("reposts缓存存储成功")
					}
				}
				session.MarkMessage(message, "")
			}
			// 清空消息列表
			messages = messages[:0]
		}
	}
}

// 启动消费者组
func (k *Kafka) Consume(ctx context.Context, consumerGroup sarama.ConsumerGroup) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := consumerGroup.Consume(ctx, k.topics, k); err != nil {
					log.Errorf("Error from consumer: %v", err)
					return
				}
			}
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
