package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/model"
	"shop_srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenerateOrderSn(userId int32) string {
	// 生成订单号的函数
	/*
		订单号的生成规则：
		年月日时分秒 + 用户id + 2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	// 获取用户的购物车列表
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse

	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		// 查询用户购物车记录
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,      // 购物车记录的唯一标识
			UserId:  shopCart.User,    // 用户id
			GoodsId: shopCart.Goods,   // 商品id
			Nums:    shopCart.Nums,    // 商品数量
			Checked: shopCart.Checked, // 是否选中
		})
	}
	return &rsp, nil
}

func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 将商品添加到购物车
	// 1. 如果购物车中原本没有这件商品则新建一个记录
	// 2. 如果这个商品之前已经添加到了购物车，则合并记录

	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		// 如果记录已经存在，则合并购物车记录，进行更新操作
		shopCart.Nums += req.Nums
	} else {
		// 如果记录不存在，则创建新的购物车记录
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 更新购物车记录的数量和选中状态
	var shopCart model.ShoppingCart

	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		// 检查购物车记录是否存在
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)

	return &emptypb.Empty{}, nil
}

func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 删除购物车记录
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		// 检查购物车记录是否存在
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// 获取订单列表
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	// 分页查询订单
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,                                      // 订单id
			UserId:  order.User,                                    // 用户id
			OrderSn: order.OrderSn,                                 // 订单号
			PayType: order.PayType,                                 // 支付方式
			Status:  order.Status,                                  // 订单状态
			Post:    order.Post,                                    // 邮费
			Total:   order.OrderMount,                              // 订单金额
			Address: order.Address,                                 // 收货地址
			Name:    order.SignerName,                              // 收货人姓名
			Mobile:  order.SingerMobile,                            // 收货人手机号
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"), // 添加时间
		})
	}
	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// 获取订单详情
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		// 检查订单是否存在
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := proto.OrderInfoResponse{
		Id:      order.ID,           // 订单id
		UserId:  order.User,         // 用户id
		OrderSn: order.OrderSn,      // 订单号
		PayType: order.PayType,      // 支付方式
		Status:  order.Status,       // 订单状态
		Post:    order.Post,         // 邮费
		Total:   order.OrderMount,   // 订单金额
		Address: order.Address,      // 收货地址
		Name:    order.SignerName,   // 收货人姓名
		Mobile:  order.SingerMobile, // 收货人手机号
	}

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		// 查询订单商品
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,      // 商品id
			GoodsName:  orderGood.GoodsName,  // 商品名称
			GoodsPrice: orderGood.GoodsPrice, // 商品价格
			GoodsImage: orderGood.GoodsImage, // 商品图片
			Nums:       orderGood.Nums,       // 商品数量
		})
	}

	return &rsp, nil
}

type OrderListener struct {
	Code        codes.Code      // 返回的状态码
	Detail      string          // 状态详情
	ID          int32           // 订单id
	OrderAmount float32         // 订单金额
	Ctx         context.Context // 上下文对象
}

// ExecuteLocalTransaction 是 OrderListener 的方法，用于执行本地事务。
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	// 解析消息中的订单信息
	var orderInfo model.OrderInfo
	err := json.Unmarshal(msg.Body, &orderInfo)
	if err != nil {
		o.Code = codes.InvalidArgument
		o.Detail = "无效的订单信息"
		return primitive.RollbackMessageState
	}

	// 查询选中结算的购物车记录
	var shopCarts []model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}

	// 提取购物车中的商品ID和数量
	var goodsIds []int32
	goodsNumsMap := make(map[int32]int32)
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 跨服务调用商品微服务查询商品信息

	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		// 计算订单总金额
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])

		// 构建订单商品详情
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		// 构建库存服务请求的商品信息
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 跨服务调用库存微服务进行库存扣减
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{OrderSn: orderInfo.OrderSn, GoodsInfo: goodsInvInfo}); err != nil {
		// 如果扣减库存失败，回滚事务
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}

	// 开启数据库事务
	tx := global.DB.Begin()

	// 保存订单信息
	orderInfo.OrderMount = orderAmount
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}

	// 关联订单商品并批量插入数据库
	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}

	// 删除购物车记录
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}

	// 发送延时消息

	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{global.RocketMQ}), producer.WithGroupName("shop_order"))
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// 将消息设置为订单超时消息，延时发送
	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(3)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	// 提交事务
	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	// 解析消息中的订单信息
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	// 检查之前的逻辑是否完成，本例中通过查询数据库中是否存在对应订单信息来判断
	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState // 之前逻辑未完成，需要回滚消息
	}

	return primitive.RollbackMessageState // 之前逻辑已完成，可以提交消息
}

func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
			3. 库存的扣减 - 访问库存服务 (跨微服务)
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/

	// 创建订单监听器
	orderListener := OrderListener{Ctx: ctx}

	// 创建 RocketMQ 事务生产者
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{global.RocketMQ}),
		producer.WithGroupName("shop_inventory"),
	)
	if err != nil {
		zap.S().Errorf("生成producer失败: %s", err.Error())
		return nil, err
	}

	// 启动生产者
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败: %s", err.Error())
		return nil, err
	}

	// 创建订单信息
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId), // 生成订单号
		Address:      req.Address,                 // 收货地址
		SignerName:   req.Name,                    // 签收人姓名
		SingerMobile: req.Mobile,                  // 签收人手机号
		Post:         req.Post,                    // 邮费
		User:         req.UserId,                  // 用户ID
	}

	// 序列化订单信息为 JSON
	jsonString, _ := json.Marshal(order)

	// 发送事务消息
	_, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}

	// 检查订单监听器返回的结果
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	// 返回订单信息响应
	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: order.OrderSn, Total: orderListener.OrderAmount}, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	// 更新订单状态

	// 先查询订单是否存在，再更新，实际上执行了两条 SQL 语句，select 和 update 语句
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// 订单超时处理方法

	for i := range msgs {
		// 解析消息中的订单信息
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())

		// 查询订单的支付状态，如果已支付则不处理，如果未支付则归还库存
		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()

			// 归还库存，可以模拟发送消息到 "order_reback" 主题中
			// 同时修改订单状态为 "TRADE_CLOSED"
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			// 创建 RocketMQ 生产者
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{global.RocketMQ}), producer.WithGroupName("inventory_reback"))
			if err != nil {
				panic("生成producer失败")
			}

			// 启动生产者
			if err = p.Start(); err != nil {
				panic("启动producer失败")
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}

			//if err = p.Shutdown(); err != nil {panic("关闭producer失败")}
			return consumer.ConsumeSuccess, nil
		}
	}

	return consumer.ConsumeSuccess, nil
}
