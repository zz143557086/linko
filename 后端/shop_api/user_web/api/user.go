package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/user_web/forms"
	"shop_api/user_web/middlewares"
	"shop_api/user_web/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/user_web/global"
	"shop_api/user_web/proto"
)

var trans ut.Translator

// 移除字段中的顶层结构体
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 处理验证器错误
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)), // 使用全局Translator翻译错误消息，并返回JSON响应
	})
	return
}

// 将gRPC的错误转换成对应的HTTP状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": string(e.Code()) + e.Message(),
				})
			}
			return
		}
	}
}

// 获取用户列表数据
func GetUserList(ctx *gin.Context) {
	// 初始化连接 gRPC 服务器
	claims, _ := ctx.Get("claims")
	CurrentUser := claims.(*models.CustomClaims)
	zap.S().Debugf("当前的用户为：" + CurrentUser.Name)
	pn := ctx.DefaultQuery("pn", "0") // 获取分页参数，默认为 0
	pnInt, _ := strconv.Atoi(pn)
	psize := ctx.DefaultQuery("psize", "5") // 获取每页大小参数，默认为 5
	pnSize, _ := strconv.Atoi(psize)

	// 调用用户服务的 GetUserList 方法，获取用户列表数据
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{Pn: uint32(pnInt), PSize: uint32(pnSize)})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx) // 处理gRPC错误，并转换为HTTP响应返回
		return
	}

	zap.S().Debug("获取用户列表页")

	result := make([]map[string]interface{}, 0) // 存储处理结果的切片

	// 将获取的用户列表数据进行处理，并添加到 result 中
	for _, value := range rsp.Data {
		data := make(map[string]interface{})
		data["Name"] = value.Name
		data["Id"] = value.Id
		data["Birthday"] = time.Unix(int64(value.BirthDay), 0).Format("2006-01-02")
		data["Gender"] = value.Gender
		data["Mobile"] = value.Mobile
		result = append(result, data)
	}

	ctx.JSON(200, result) // 将 result 以 JSON 格式作为响应发送给客户端
}

// 密码登录处理函数
func PassWordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	// 表单验证
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		fmt.Println(err.Error(), "错误")
		HandleValidatorError(c, err) // 处理验证器错误
		return
	}

	if store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 查询用户手机号
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		// 判断错误类型
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				// 用户不存在
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				// 登录失败，系统错误
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败,系统错误",
				})
			}
		}
	} else {
		// 只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password: passwordLoginForm.PassWord,
			Name:     rsp.Name,
		}); pasErr != nil {
			// 登录失败，系统错误
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败,系统错误",
			})
		} else {
			if passRsp.Success {
				// 登录成功
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					Name:        rsp.Name,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "linko",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"name":       rsp.Name,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				// 登录失败，密码错误
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败,密码错误",
				})
			}
		}
	}
}
func Register(c *gin.Context) {
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 调用用户服务的 CreateUserInfo 方法创建用户
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Name:     registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 生成 JWT token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		Name:        user.Name,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	// 返回注册成功的响应给客户端
	c.JSON(200, gin.H{
		"id":         user.Id,
		"nick_name":  user.Name,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
