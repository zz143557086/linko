package forms

// PassWordLoginForm 密码登录表单
type PassWordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`           // 手机号码，必填，格式验证为手机号格式
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"` // 密码，必填，长度限制为3-20
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`    // 验证码，必填，长度限制为5
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`          // 验证码ID，必填
}

// RegisterForm 注册表单
type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`           // 手机号码，必填，格式验证为手机号格式
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"` // 密码，必填，长度限制为3-20
}

// UpdateUserForm 更新用户表单
type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`                // 姓名，必填，长度限制为3-10
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`       // 性别，必填，只能是female或male
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"` // 生日，必填，需要符合指定的日期格式（2006-01-02）
}
