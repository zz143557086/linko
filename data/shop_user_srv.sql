/*
 Navicat Premium Data Transfer

 Source Server         : 我的虚拟机
 Source Server Type    : MySQL
 Source Server Version : 50736 (5.7.36)
 Source Host           : 192.168.2.106:3306
 Source Schema         : shop_user_srv

 Target Server Type    : MySQL
 Target Server Version : 50736 (5.7.36)
 File Encoding         : 65001

 Date: 15/08/2023 21:45:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `add_time` datetime(3) NULL DEFAULT NULL COMMENT '新增时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '跟新时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint(1) NULL DEFAULT NULL COMMENT '是否删除',
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '手机号码',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '姓名',
  `birthday` datetime NULL DEFAULT NULL COMMENT '生日',
  `gender` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'male' COMMENT 'female表示女, male表示男',
  `role` int(11) NULL DEFAULT 1 COMMENT '1表示普通用户, 2表示管理员',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `mobile`(`mobile`) USING BTREE,
  INDEX `idx_mobile`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '2023-07-24 22:27:14.991', '2023-07-24 22:27:14.991', NULL, 0, '15520050430', 'e10adc3949ba59abbe56e057f20f883e', 'lin0', '2023-07-25 14:22:15', 'male', 2);
INSERT INTO `user` VALUES (2, '2023-07-24 22:27:15.041', '2023-07-24 22:27:15.041', NULL, 0, '15520050431', 'e10adc3949ba59abbe56e057f20f883e', 'lin1', '2023-07-25 14:22:13', 'male', 2);
INSERT INTO `user` VALUES (3, '2023-07-24 22:27:15.075', '2023-07-24 22:27:15.075', NULL, 0, '15520050432', 'e10adc3949ba59abbe56e057f20f883e', 'lin2', '2023-07-25 14:22:11', 'male', 2);
INSERT INTO `user` VALUES (4, '2023-07-24 22:27:15.109', '2023-07-24 22:27:15.109', NULL, 0, '15520050433', 'e10adc3949ba59abbe56e057f20f883e', 'lin3', '2023-07-25 14:22:09', 'male', 2);
INSERT INTO `user` VALUES (5, '2023-07-24 22:27:15.137', '2023-07-24 22:27:15.137', NULL, 0, '15520050434', 'e10adc3949ba59abbe56e057f20f883e', 'lin4', '2023-07-25 14:22:18', 'male', 2);
INSERT INTO `user` VALUES (6, '2023-07-24 22:27:15.156', '2023-07-24 22:27:15.156', NULL, 0, '15520050435', 'e10adc3949ba59abbe56e057f20f883e', 'lin5', '2023-07-25 14:22:07', 'male', 1);
INSERT INTO `user` VALUES (7, '2023-07-24 22:27:15.178', '2023-07-24 22:27:15.178', NULL, 0, '15520050436', 'e10adc3949ba59abbe56e057f20f883e', 'lin6', '2023-07-25 14:22:03', 'male', 1);
INSERT INTO `user` VALUES (8, '2023-07-24 22:27:15.193', '2023-07-24 22:27:15.193', NULL, 0, '15520050437', 'e10adc3949ba59abbe56e057f20f883e', 'lin7', '2023-07-25 14:22:05', 'male', 1);
INSERT INTO `user` VALUES (9, '2023-07-24 22:27:15.218', '2023-07-24 22:27:15.218', NULL, 0, '15520050438', 'e10adc3949ba59abbe56e057f20f883e', 'lin8', '2023-07-25 14:22:20', 'male', 1);
INSERT INTO `user` VALUES (10, '2023-07-24 22:27:15.244', '2023-07-24 22:27:15.244', NULL, 0, '15520050439', 'e10adc3949ba59abbe56e057f20f883e', 'lin9', '2023-07-25 14:21:58', 'male', 1);
INSERT INTO `user` VALUES (11, '2023-07-25 09:20:41.354', '2023-07-25 09:20:41.354', NULL, 0, '11111111110', 'e10adc3949ba59abbe56e057f20f883e', 'joker0', '2023-07-25 14:22:01', 'male', 1);
INSERT INTO `user` VALUES (12, '2023-07-25 09:20:41.387', '2023-07-25 09:20:41.387', NULL, 0, '11111111111', 'e10adc3949ba59abbe56e057f20f883e', 'joker1', '2023-07-25 14:21:56', 'male', 1);
INSERT INTO `user` VALUES (13, '2023-07-25 09:20:41.432', '2023-07-25 09:20:41.432', NULL, 0, '11111111112', 'e10adc3949ba59abbe56e057f20f883e', 'joker2', '2023-07-25 14:22:22', 'male', 1);
INSERT INTO `user` VALUES (14, '2023-07-25 09:20:41.477', '2023-07-25 09:20:41.477', NULL, 0, '11111111113', 'e10adc3949ba59abbe56e057f20f883e', 'joker3', '2023-07-25 14:21:54', 'male', 1);
INSERT INTO `user` VALUES (15, '2023-07-25 09:20:41.503', '2023-07-25 09:20:41.503', NULL, 0, '11111111114', 'e10adc3949ba59abbe56e057f20f883e', 'joker4', '2023-07-25 14:22:25', 'male', 1);
INSERT INTO `user` VALUES (16, '2023-07-25 09:20:41.531', '2023-07-25 09:20:41.531', NULL, 0, '11111111115', 'e10adc3949ba59abbe56e057f20f883e', 'joker5', '2023-07-25 14:21:45', 'male', 1);
INSERT INTO `user` VALUES (17, '2023-07-25 09:20:41.568', '2023-07-25 09:20:41.568', NULL, 0, '11111111116', 'e10adc3949ba59abbe56e057f20f883e', 'joker6', '2023-07-25 14:21:49', 'male', 1);
INSERT INTO `user` VALUES (18, '2023-07-25 09:20:41.593', '2023-07-25 09:20:41.593', NULL, 0, '11111111117', 'e10adc3949ba59abbe56e057f20f883e', 'joker7', '2023-07-25 14:21:51', 'male', 1);
INSERT INTO `user` VALUES (19, '2023-07-25 09:20:41.621', '2023-07-25 09:20:41.621', NULL, 0, '11111111118', 'e10adc3949ba59abbe56e057f20f883e', 'joker8', '2023-07-25 14:22:27', 'male', 1);
INSERT INTO `user` VALUES (20, '2023-07-25 09:20:41.641', '2023-07-25 09:20:41.641', NULL, 0, '11111111119', 'e10adc3949ba59abbe56e057f20f883e', 'joker9', '2023-07-25 14:22:30', 'male', 1);
INSERT INTO `user` VALUES (21, '2023-07-27 10:25:48.147', '2023-07-27 10:25:48.147', NULL, 0, '17723575440', 'e10adc3949ba59abbe56e057f20f883e', '17723575440', NULL, 'male', 1);

SET FOREIGN_KEY_CHECKS = 1;
