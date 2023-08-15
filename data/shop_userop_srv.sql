/*
 Navicat Premium Data Transfer

 Source Server         : 我的虚拟机
 Source Server Type    : MySQL
 Source Server Version : 50736 (5.7.36)
 Source Host           : 192.168.2.106:3306
 Source Schema         : shop_userop_srv

 Target Server Type    : MySQL
 Target Server Version : 50736 (5.7.36)
 File Encoding         : 65001

 Date: 15/08/2023 21:45:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for address
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '\'用户ID\'',
  `province` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'省份\'',
  `city` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'城市\'',
  `district` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'区/县\'',
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'详细地址\'',
  `signer_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'收件人姓名\'',
  `signer_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'收件人手机号\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_address_user`(`user`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of address
-- ----------------------------

-- ----------------------------
-- Table structure for leavingmessages
-- ----------------------------
DROP TABLE IF EXISTS `leavingmessages`;
CREATE TABLE `leavingmessages`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '\'用户ID\'',
  `message_type` int(11) NULL DEFAULT NULL COMMENT '\'留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)\'',
  `subject` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'主题\'',
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '\'留言内容\'',
  `file` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '\'文件路径\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_leavingmessages_user`(`user`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of leavingmessages
-- ----------------------------
INSERT INTO `leavingmessages` VALUES (1, '2023-08-14 14:47:47.924', '2023-08-14 14:47:47.924', NULL, 0, 1, 1, 'officia ut deserunt sint eiusmod', 'consectetur eiusmod qui', 'wais://tvysbqlifd.hr/iqlqk');

-- ----------------------------
-- Table structure for userfav
-- ----------------------------
DROP TABLE IF EXISTS `userfav`;
CREATE TABLE `userfav`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '\'用户ID\'',
  `goods` int(11) NULL DEFAULT NULL COMMENT '\'商品ID\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_user_goods`(`user`, `goods`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of userfav
-- ----------------------------
INSERT INTO `userfav` VALUES (1, '2023-08-14 15:02:14.186', '2023-08-14 15:02:14.186', NULL, 0, 1, 424);

SET FOREIGN_KEY_CHECKS = 1;
