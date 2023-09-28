/*
 Navicat Premium Data Transfer

 Source Server         : 本地库
 Source Server Type    : MySQL
 Source Server Version : 80100
 Source Host           : localhost:3306
 Source Schema         : n9e_v6

 Target Server Type    : MySQL
 Target Server Version : 80100
 File Encoding         : 65001

 Date: 28/09/2023 14:19:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for device_type
-- ----------------------------
DROP TABLE IF EXISTS `device_type`;
CREATE TABLE `device_type`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '名称',
  `TYPES` int(0) NOT NULL COMMENT '类别(1:设备类型;2:备件设备类型)',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '设备类型' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of device_type
-- ----------------------------
INSERT INTO `device_type` VALUES (1, 'ARM服务器', 1, 'root', 1692858966, '', 1692858966);
INSERT INTO `device_type` VALUES (2, '备份设备', 1, 'root', 1692858990, '', 1692858990);
INSERT INTO `device_type` VALUES (3, '环境动力', 1, 'root', 1692859007, '', 1692859007);
INSERT INTO `device_type` VALUES (4, '超融合', 1, 'root', 1692859014, '', 1692859014);
INSERT INTO `device_type` VALUES (5, '工控机', 1, 'root', 1692859032, '', 1692859032);
INSERT INTO `device_type` VALUES (6, '大型机', 1, 'root', 1692859038, '', 1692859038);
INSERT INTO `device_type` VALUES (7, '小型机', 1, 'root', 1692859044, '', 1692859044);
INSERT INTO `device_type` VALUES (8, 'PC机', 1, 'root', 1692859058, '', 1692859058);
INSERT INTO `device_type` VALUES (9, '网络设备', 1, 'root', 1692859068, '', 1692859068);
INSERT INTO `device_type` VALUES (10, '其他', 1, 'root', 1692859075, '', 1692859075);
INSERT INTO `device_type` VALUES (11, '安全设备', 1, 'root', 1692859092, '', 1692859092);
INSERT INTO `device_type` VALUES (12, 'X86服务器', 1, 'root', 1692859101, '', 1692859101);
INSERT INTO `device_type` VALUES (13, '存储', 1, 'root', 1692859114, '', 1692859114);
INSERT INTO `device_type` VALUES (14, '终端设备', 1, 'root', 1692859129, '', 1692859129);

SET FOREIGN_KEY_CHECKS = 1;
