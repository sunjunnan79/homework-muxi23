package jwt

import (
	. "Insomnia/app/infrastructure/config"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// SignToken 生成 JWT 令牌
func SignToken(Uuid string) (string, error) {
	config := LoadConfig() // 加载应用程序配置信息
	secretKey := config.JWT.JWTSecretKey
	issuer := config.JWT.Issuer

	// 设置过期时间为当前时间加一天
	expirationTime := time.Now().Add(time.Hour * 24)

	// 创建声明
	claims := jwt.MapClaims{
		"exp":  expirationTime.Unix(),
		"iss":  issuer,
		"iat":  time.Now().Unix(),
		"Uuid": Uuid,
		"nbf":  time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err // 如果签名过程中出现错误，返回错误信息
	}

	return signedToken, nil // 返回签名后的 JWT 令牌
}

// Parse 解析JWT令牌并返回声明
func Parse(tokenString string) (jwt.MapClaims, error) {
	config := LoadConfig() // 加载应用程序配置信息
	secretKey := config.JWT.JWTSecretKey
	// 解析JWT令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 提取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	return claims, nil
}
