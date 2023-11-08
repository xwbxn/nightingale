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

 Date: 26/10/2023 09:14:55
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dict_type
-- ----------------------------
DROP TABLE IF EXISTS `dict_type`;
CREATE TABLE `dict_type`  (
  `ID` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `TYPE_CODE` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典编码',
  `DICT_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典名称',
  `IS_VISIBLE` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '是否可见',
  `REMARK` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` bigint NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` bigint NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典类别表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of dict_type
-- ----------------------------
INSERT INTO `dict_type` VALUES (27, 'device_status', '设备状态', 'YES', '', 'root', 1692846004, '', 1692846004, NULL);
INSERT INTO `dict_type` VALUES (28, 'operate_system', '操作系统', 'YES', '', 'root', 1692926896, '', 1692926896, NULL);
INSERT INTO `dict_type` VALUES (29, 'maintenance_type', '维保类型', 'YES', '', 'root', 1693361092, '', 1693361092, NULL);
INSERT INTO `dict_type` VALUES (30, 'maintenance_service', '维保服务', 'YES', '', 'root', 1694076198, '', 1694076198, NULL);
INSERT INTO `dict_type` VALUES (31, 'alert_event', '变更事项', 'YES', '', 'root', 1694411172, '', 1694411172, NULL);
INSERT INTO `dict_type` VALUES (32, 'service_level', '服务级别', 'YES', '', 'root', 1694569480, '', 1694569480, NULL);
INSERT INTO `dict_type` VALUES (33, 'producer-type', '厂商类型', 'NO', '', 'root', 1695173905, '', 1695173905, NULL);
INSERT INTO `dict_type` VALUES (34, 'spare-base-data', '备件基础数据', 'YES', '', 'root', 1695174448, '', 1695174448, NULL);
INSERT INTO `dict_type` VALUES (35, 'sub_type', '子类型', 'YES', '', 'root', 1695285102, '', 1695285102, NULL);
INSERT INTO `dict_type` VALUES (36, 'outline_structure', '外形结构', 'YES', '', 'root', 1695285171, '', 1695285171, NULL);
INSERT INTO `dict_type` VALUES (37, 'out_band_version', '带外版本', 'YES', '', 'root', 1695285725, '', 1695285725, NULL);
INSERT INTO `dict_type` VALUES (38, 'basic_expansion', '基本信息', 'YES', '', 'root', 1695366230, '', 1695366230, NULL);
INSERT INTO `dict_type` VALUES (39, 'asset_expansion_fields', '扩展字段', 'NO', '', 'root', 1695604819, 'root', 1695604840, NULL);
INSERT INTO `dict_type` VALUES (41, 'link_method', '连接方式', 'YES', '', 'root', 1695696617, '', 1695696617, NULL);
INSERT INTO `dict_type` VALUES (42, 'asset_filter_status', '设备资产过滤器状态', 'YES', '', 'root', 1695712501, '', 1695712501, NULL);
INSERT INTO `dict_type` VALUES (43, 'offline_release_status', '下线释放资源', 'NO', '', 'root', 1695806497, '', 1695806497, NULL);
INSERT INTO `dict_type` VALUES (44, 'scrap_query_filter', '报废设备过滤器', 'YES', '', 'root', 1696749675, '', 1696749675, NULL);
INSERT INTO `dict_type` VALUES (45, 'system_runing', '系统运行时间', 'YES', '', 'root', 1697680946, '', 1697680946, '2023-10-19 10:12:39.258');
INSERT INTO `dict_type` VALUES (46, 'system_initialize', '系统初始化参数', 'YES', '', 'root', 1697681088, '', 1697681088, NULL);
INSERT INTO `dict_type` VALUES (47, 'asset_ext_fields', '资产扩展属性', 'YES', '', 'root', 1697780263, '', 1697780263, NULL);
INSERT INTO `dict_type` VALUES (48, 'cpu_specifications', 'CPU规格', 'YES', '', 'root', 1697785560, '', 1697785560, NULL);
INSERT INTO `dict_type` VALUES (49, 'memory_specifications', '内存规格', 'YES', '', 'root', 1697785682, '', 1697785682, NULL);
INSERT INTO `dict_type` VALUES (50, 'producer', '厂商', 'YES', '', 'root', 1697785758, '', 1697785758, NULL);
INSERT INTO `dict_type` VALUES (51, 'asset_status', '资产状态', 'YES', '', 'root', 1697785853, '', 1697785853, NULL);

SET FOREIGN_KEY_CHECKS = 1;
