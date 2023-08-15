/*
 Navicat Premium Data Transfer

 Source Server         : 我的虚拟机
 Source Server Type    : MySQL
 Source Server Version : 50736 (5.7.36)
 Source Host           : 192.168.2.106:3306
 Source Schema         : shop_order_srv

 Target Server Type    : MySQL
 Target Server Version : 50736 (5.7.36)
 File Encoding         : 65001

 Date: 15/08/2023 21:45:08
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ordergoods
-- ----------------------------
DROP TABLE IF EXISTS `ordergoods`;
CREATE TABLE `ordergoods`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `order` int(11) NULL DEFAULT NULL COMMENT '订单ID',
  `goods` int(11) NULL DEFAULT NULL COMMENT '商品ID',
  `goods_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '商品名称',
  `goods_image` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '商品图片',
  `goods_price` float NULL DEFAULT NULL COMMENT '商品价格',
  `nums` int(11) NULL DEFAULT NULL COMMENT '商品数量',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ordergoods_order`(`order`) USING BTREE,
  INDEX `idx_ordergoods_goods`(`goods`) USING BTREE,
  INDEX `idx_ordergoods_goods_name`(`goods_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ordergoods
-- ----------------------------
INSERT INTO `ordergoods` VALUES (8, '2023-08-13 16:18:40.172', '2023-08-13 16:18:40.172', NULL, 0, 8, 423, '越南进口红心火龙果 4个装 红肉中果 单果约330-420g 新鲜水果', 'https://img12.360buyimg.com/n5/jfs/t26575/223/1308050262/395693/66e2d658/5bc69bcfN8030a03e.jpg', 27.9, 1);
INSERT INTO `ordergoods` VALUES (10, '2023-08-14 09:25:13.990', '2023-08-14 09:25:13.990', NULL, 0, 10, 423, '越南进口红心火龙果 4个装 红肉中果 单果约330-420g 新鲜水果', 'https://img12.360buyimg.com/n5/jfs/t26575/223/1308050262/395693/66e2d658/5bc69bcfN8030a03e.jpg', 27.9, 3);

-- ----------------------------
-- Table structure for orderinfo
-- ----------------------------
DROP TABLE IF EXISTS `orderinfo`;
CREATE TABLE `orderinfo`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '用户ID',
  `order_sn` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '订单号',
  `pay_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'alipay(支付宝)， wechat(微信)\'',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)\'',
  `trade_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '交易号',
  `order_mount` float NULL DEFAULT NULL COMMENT '订单金额',
  `pay_time` datetime NULL DEFAULT NULL COMMENT '支付时间',
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址',
  `signer_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '签收人姓名',
  `singer_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '签收人手机号',
  `post` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '邮编',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_orderinfo_user`(`user`) USING BTREE,
  INDEX `idx_orderinfo_order_sn`(`order_sn`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of orderinfo
-- ----------------------------
INSERT INTO `orderinfo` VALUES (8, '2023-08-13 16:18:40.154', '2023-08-13 16:18:40.154', NULL, 0, 2, '2023813161817493200238', '', '', '', 27.9, NULL, '北京市', 'lin1', '15520050431', '请尽快发货');
INSERT INTO `orderinfo` VALUES (10, '2023-08-14 09:25:13.947', '2023-08-14 09:25:13.947', NULL, 0, 1, '2023814925798136100170', '', '', '', 83.7, NULL, '北京市', 'lin0', '15520050430', '请尽快发货');

-- ----------------------------
-- Table structure for shoppingcart
-- ----------------------------
DROP TABLE IF EXISTS `shoppingcart`;
CREATE TABLE `shoppingcart`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '用户ID',
  `goods` int(11) NULL DEFAULT NULL COMMENT '商品ID',
  `nums` int(11) NULL DEFAULT NULL COMMENT '商品数量',
  `checked` tinyint(1) NULL DEFAULT NULL COMMENT '是否选中',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_shoppingcart_user`(`user`) USING BTREE,
  INDEX `idx_shoppingcart_goods`(`goods`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shoppingcart
-- ----------------------------
INSERT INTO `shoppingcart` VALUES (3, '2023-08-11 12:11:30.426', '2023-08-11 12:11:30.426', NULL, 0, 1, 420, 1, 0);
INSERT INTO `shoppingcart` VALUES (4, '2023-08-11 12:12:09.622', '2023-08-11 15:27:06.925', '2023-08-14 09:25:14.000', 0, 1, 423, 3, 1);
INSERT INTO `shoppingcart` VALUES (5, '2023-08-13 12:05:53.468', '2023-08-13 12:05:53.468', '2023-08-13 16:18:40.278', 0, 2, 423, 1, 1);
INSERT INTO `shoppingcart` VALUES (6, '2023-08-14 10:10:05.291', '2023-08-14 10:59:49.752', '2023-08-14 11:00:10.838', 0, 1, 425, 6, 0);
INSERT INTO `shoppingcart` VALUES (7, '2023-08-14 10:11:11.591', '2023-08-14 10:11:11.591', '2023-08-14 11:00:13.037', 0, 1, 455, 2, 0);

SET FOREIGN_KEY_CHECKS = 1;
