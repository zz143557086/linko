package global

import (
	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator
	Tags  = []string{"linko", "lin", "oss_web"}
)

const (
	ClientPort    = 8002
	ServerPort    = 8001
	Host          = "127.0.0.1"
	ServerHost    = "127.0.0.1"
	ClientHost    = "127.0.0.1"
	ServerName    = "oss_srv"
	ClientName    = "oss_web"
	ConsulHost    = "192.168.2.106"
	ConsulPort    = 8500
	OssUploadDir  = "linko/"
	OssApiSecrect = "G37bhtUruBt5YE2TeP7zvSiFixZNbo"
	OssApiKeyId   = "LTAI5tDi2S4Y1ffQrPrK872s"
	CallBackUrl   = "http://88.88.88.88:8888"
	OssSigningKey = 3000
	Key           = "ElwvtT$G@ceMdGE@rsBAZc9vOFs9zqO4"
	OssHost       = "http://linko-shop.oss-cn-hangzhou.aliyuncs.com"
)
