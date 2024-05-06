package Email

import (
	"Insomnia/app/infrastructure/redis"
	. "Insomnia/app/utility/tool"
	"context"
	"errors"
	"time"
)

type CheckEmail struct {
	Email            string `gorm:"size:255;not null;unique"`
	VerificationCode string `gorm:"size:255;not null"`
}

// CreateRedis 创建一个Redis
func (email CheckEmail) CreateRedis() (err error) {
	//加密密码
	email.VerificationCode = Encrypt(email.VerificationCode)
	ctx := context.Background()
	key := "verification_code:" + email.Email
	//限制时间为5分钟
	expiration := 5 * time.Minute
	//设施一个redis
	err = redis.Rdb.Set(ctx, key, email.VerificationCode, expiration).Err()
	if err != nil {
		return
	}
	return
}

// CheckVerificationCode 方法检查验证码是否正确
func (email CheckEmail) CheckVerificationCode() (err error) {
	email.VerificationCode = Encrypt(email.VerificationCode)
	// 在 Redis 中检查验证码
	ctx := context.Background()
	key := "verification_code:" + email.Email
	storedCode, err := redis.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if storedCode != email.VerificationCode {
		return errors.New("无效的验证码")
	}
	return
}
