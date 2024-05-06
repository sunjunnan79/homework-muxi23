package redis

import (
	. "Insomnia/app/infrastructure/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os/exec"
	"time"
)

var Rdb *redis.Client

func init() {
	// 检查 Redis 服务器是否已经运行
	redisRunning, err := isRedisRunning()
	if err != nil {
		fmt.Println("Failed to check Redis status:", err)
		return
	}
	// 如果 Redis 服务器未运行，则尝试启动 Redis 服务器
	if !redisRunning {
		err := startRedisServer()
		if err != nil {
			fmt.Println("Failed to start Redis server:", err)
			return
		}
		fmt.Println("Redis server started.")
		config := LoadConfig()
		//初始化Redis客户端
		Rdb = redis.NewClient(&redis.Options{
			Addr:     config.Re.Address,
			Password: config.Re.Password,
			DB:       config.Re.Database,
		})
	}

}

// 检查 Redis 服务器是否已经运行
func isRedisRunning() (bool, error) {
	ctx := context.Background()
	config := LoadConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Re.Address,
		Password: config.Re.Password,
		DB:       config.Re.Database,
	})

	_, err := client.Ping(ctx).Result()
	if err == nil {
		Rdb = client
		return true, nil
	}
	defer client.Close()
	if err.Error() == "redis: can't connect to the server" {
		return false, err
	}

	return false, nil
}

// 启动 Redis 服务器
func startRedisServer() error {
	cmd := exec.Command("redis-server")
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// SetResp 设置redis里面的结果
func SetResp(owner string, value any) error {
	return Rdb.Set(context.Background(), owner, value, 1*time.Hour).Err()
}

// SetNXResp 如果没设置过的话就设置
func SetNXResp(owner string, value any) error {
	return Rdb.SetNX(context.Background(), owner, value, 1*time.Hour).Err()
}

// ExistResp 检查是否已经存储
func ExistResp(owner string) (exist int64, err error) {
	exist, err = Rdb.Exists(context.Background(), owner).Result()
	return
}

// DelResp 删除对应的缓存
func DelResp(owner string) error {
	return Rdb.Del(context.Background(), owner).Err()
}

func GetResp(owner string) ([]byte, error) {
	get := Rdb.Get(context.Background(), owner)
	if get.Err() != nil {
		return nil, get.Err()
	}
	bytes, err := get.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
