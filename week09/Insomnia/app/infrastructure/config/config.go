package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type App struct {
	Address string
	Static  string
	Log     string
}

type Database struct {
	Driver   string
	Server   string
	Port     string
	Database string
	User     string
	Password string
	Config   string
}

type Redis struct {
	Address  string
	Database int
	Password string
}

type Email struct {
	UserName string
	Sender   string
	Password string
	Smtp     string
}

type Configuration struct {
	App   App
	Db    Database
	Re    Redis
	Email Email
	Oss   Oss
	JWT   JWT
}

type Oss struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

type JWT struct {
	JWTSecretKey string
	Issuer       string
}

var config *Configuration
var once sync.Once

func LoadConfig() *Configuration {
	once.Do(func() {
		file, err := os.Open("config.json")
		if err != nil {
			log.Fatalln("无法打开 config 文件", err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		config = &Configuration{}
		err = decoder.Decode(config)
		if err != nil {
			log.Fatalln("无法从文件获取 configuration", err)
		}
	})
	return config
}
