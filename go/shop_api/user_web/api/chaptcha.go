package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

// GetCaptcha 获取验证码接口
func GetCaptcha(ctx *gin.Context) {
	// 创建一个基于数字的验证码驱动器，参数依次为图片宽度、高度、验证码位数、噪音强度、干扰线数
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码，返回验证码的ID、图片的base64编码以及可能的错误
	id, b64s, err := cp.Generate()
	if err != nil {
		// 如果生成验证码时发生错误，记录错误日志，并返回错误响应
		zap.S().Errorf("生成验证码错误: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}

	// 返回成功响应，包含验证码ID和图片的base64编码
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
