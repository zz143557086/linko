package hander

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/proto"
)

type UserServer struct {
}

// Paginate 根据提供的页面和页面大小参数，返回用于设置gorm.DB查询的函数
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 处理页码
		if page <= 0 {
			page = 1
		}

		// 处理页面大小
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		// 计算偏移量
		offset := (page - 1) * pageSize

		// 设置偏移量和限制数量
		return db.Offset(offset).Limit(pageSize)
	}
}

// ModelToResponse 将model.User转换为proto.UserInfoResponse
func ModelToResponse(user model.User) proto.UserInfoResponse {
	UserInfoRsp := proto.UserInfoResponse{
		Id:     user.ID,
		Name:   user.Name,
		Gender: user.Gender,
		Role:   user.Role,
		Mobile: user.Mobile,
	}
	if user.Birthday != nil {
		UserInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return UserInfoRsp
}

// genMd5 对字符串进行MD5加密，返回生成的MD5哈希值的十六进制表示
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

// GetUserList 根据提供的页面信息获取用户列表
func (s UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 查询所有用户
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	// 分页查询用户
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	// 将查询到的用户转换为响应类型
	for _, user := range users {
		UserInfoRespon := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &UserInfoRespon)
	}

	return &rsp, nil
}

// GetUserByMobile 根据提供的手机号请求获取用户信息
func (s UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 获取手机号
	mobile := req.Mobile

	// 查询用户
	var user model.User
	result := global.DB.Where("Mobile = ?", mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// 将查询到的用户转换为响应类型
	rsp := ModelToResponse(user)
	return &rsp, nil
}

// GetUserById 根据提供的用户ID请求获取用户信息
func (s UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// 获取用户ID
	id := req.Id

	// 查询用户
	var user model.User
	result := global.DB.First(&user, id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// 将查询到的用户转换为响应类型
	rsp := ModelToResponse(user)
	return &rsp, nil
}

// CreateUser 创建用户
func (s UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 查询手机号是否已存在
	var user model.User
	result := global.DB.Where("Mobile = ?", req.Mobile).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	// 创建用户
	user.Name = req.Name
	user.Password = genMd5(req.PassWord)
	user.Mobile = req.Mobile
	result = global.DB.Create(&user)
	rsp := ModelToResponse(user)
	return &rsp, nil
}

// UpdateUser 更新用户信息
func (s UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	// 查询用户是否存在
	var user model.User
	result := global.DB.Where("Id = ?", req.Id).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	// 更新用户信息
	Birthday := time.Unix(int64(req.BirthDay), 0)
	user.Name = req.Name
	user.Gender = req.Gender
	user.Birthday = &Birthday
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

// CheckPassWord 校验密码
func (s UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	// 查询用户
	var user model.User
	result := global.DB.Where("Name = ?", req.Name).First(&user)
	password := genMd5(req.Password)
	if result.Error != nil {
		return nil, result.Error
	}

	// 校验密码
	if password != user.Password {
		return &proto.CheckResponse{Success: false}, nil
	}

	return &proto.CheckResponse{Success: true}, nil
}
