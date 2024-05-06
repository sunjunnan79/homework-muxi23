package tube

import (
	. "Insomnia/app/infrastructure/config"
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"path"
	"time"
)

type Qiniu struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

var Q Qiniu

// LoadQiniu 获取七牛云的配置信息
func LoadQiniu() {
	//获取全局配置
	c := LoadConfig()
	Q = Qiniu{
		AccessKey: c.Oss.AccessKey,
		SecretKey: c.Oss.SecretKey,
		Bucket:    c.Oss.Bucket,
		Domain:    c.Oss.Domain,
	}
}

func UploadFileToQiniu(localFilePath string) (string, error) {
	mac := qbox.NewMac(Q.AccessKey, Q.SecretKey)
	//设置上传的形式与格式
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	//将cfg解析成对应格式上传
	uploader := storage.NewFormUploader(&cfg)

	//put政策的规定
	putPolicy := storage.PutPolicy{
		Scope: Q.Bucket,
	}

	//将token解析成对应格式上传
	token := putPolicy.UploadToken(mac)

	//上传回复内容(看不懂)
	ret := storage.PutRet{}

	//读取本地文件?
	remoteFileName := "captcha/" + time.Now().String() + path.Base(localFilePath)

	//上传
	err := uploader.PutFile(context.Background(), &ret, token, remoteFileName, localFilePath, nil)
	if err != nil {
		return "", err
	}

	//返回的就是图片的网址
	return Q.Domain + "/" + ret.Key, nil

}

func GetQNToken() string {

	mac := qbox.NewMac(Q.AccessKey, Q.SecretKey)

	//put政策的规定
	putPolicy := storage.PutPolicy{
		Scope: Q.Bucket,
	}

	//将token转化为对应格式
	QNToken := putPolicy.UploadToken(mac)
	return QNToken
}
