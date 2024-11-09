package tool

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"log"
	rand1 "math/rand"
	"strconv"
	"time"
)

// CreateUuid 创建uuid
func CreateUuid() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		log.Fatal("Cannot generate Uuid", err)
	}

	//看不懂的随机生成部分(看懂了是位运算)
	//0x40 是RFC 4122 中保留的变体
	u[8] = (u[8] | 0x40) & 0x7F
	//将time_hi_and_version字段的四个最高有效位(第12到15位),设置为4位版本号
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// Encrypt 使用SHA-1对明文密码进行哈希加密
func Encrypt(plaintext string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
}

func GetRandom() string {
	// 设置随机数种子
	rand1.Seed(time.Now().UnixNano())
	// 生成一个四位的随机数
	random := rand1.Intn(9000) + 1000
	//转化成string类型
	return strconv.Itoa(random)
}
