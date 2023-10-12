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

 Date: 12/10/2023 16:31:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dict_data
-- ----------------------------
DROP TABLE IF EXISTS `dict_data`;
CREATE TABLE `dict_data`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `TYPE_CODE` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典编码',
  `DICT_KEY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典标签',
  `DICT_VALUE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典键值',
  `SN` int(0) DEFAULT NULL COMMENT '序号',
  `REMARK` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 129 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '字典数据表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of dict_data
-- ----------------------------
INSERT INTO `dict_data` VALUES (31, 'ceshi', '777', '8888', NULL, '', 'root', 1692773818, '', 1692773818, NULL);
INSERT INTO `dict_data` VALUES (54, 'device_status', '未纳管V', '未纳管', NULL, '', 'root', 1692858365, '', 1692858365, NULL);
INSERT INTO `dict_data` VALUES (55, 'device_status', '已纳管V', '已纳管', NULL, '', 'root', 1692858365, '', 1692858365, NULL);
INSERT INTO `dict_data` VALUES (56, 'device_status', '未上架V', '未上架', NULL, '', 'root', 1692858365, '', 1692858365, NULL);
INSERT INTO `dict_data` VALUES (57, 'device_status', '已上架V', '已上架', NULL, '', 'root', 1692858365, '', 1692858365, NULL);
INSERT INTO `dict_data` VALUES (58, 'device_status', '关机V', '关机', NULL, '', 'root', 1692858365, '', 1692858365, NULL);
INSERT INTO `dict_data` VALUES (59, 'operate_system', '1', 'Windows', NULL, '', 'root', 1692926922, '', 1692926922, NULL);
INSERT INTO `dict_data` VALUES (60, 'operate_system', '2', 'Centos', NULL, '', 'root', 1692926922, '', 1692926922, NULL);
INSERT INTO `dict_data` VALUES (69, 'alert_event', 'parts_alert', '部件变更', NULL, '', 'root', 1694411498, '', 1694411498, NULL);
INSERT INTO `dict_data` VALUES (70, 'alert_event', 'disk_capacity_alert', '磁盘容量变更', NULL, '', 'root', 1694411498, '', 1694411498, NULL);
INSERT INTO `dict_data` VALUES (71, 'service_level', '白金', '白金', NULL, '', 'root', 1694569589, '', 1694569589, NULL);
INSERT INTO `dict_data` VALUES (72, 'service_level', '金牌', '金牌', NULL, '', 'root', 1694569589, '', 1694569589, NULL);
INSERT INTO `dict_data` VALUES (75, 'maintenance_service', '24小时X7天', '24小时X7天', NULL, '整机', 'root', 1694675057, '', 1694675057, NULL);
INSERT INTO `dict_data` VALUES (76, 'maintenance_service', '24小时', '24小时', NULL, '', 'root', 1694675057, '', 1694675057, NULL);
INSERT INTO `dict_data` VALUES (81, 'spare-base-data', 'spare_party_base', '备件基础数据', NULL, '', 'root', 1695265042, '', 1695265042, NULL);
INSERT INTO `dict_data` VALUES (82, 'spare-base-data', 'party_type', '部件类型', NULL, '', 'root', 1695265042, '', 1695265042, NULL);
INSERT INTO `dict_data` VALUES (83, 'spare-base-data', 'device_type', '设备类型', NULL, '', 'root', 1695265042, '', 1695265042, NULL);
INSERT INTO `dict_data` VALUES (84, 'spare-base-data', 'inventory_alert_set', '库存预警设置', NULL, '', 'root', 1695265042, '', 1695265042, NULL);
INSERT INTO `dict_data` VALUES (85, 'spare-base-data', 'warehouse_information', '库房信息', NULL, '', 'root', 1695265042, '', 1695265042, NULL);
INSERT INTO `dict_data` VALUES (86, 'sub_type', '刀箱', '刀箱', NULL, '', 'root', 1695285136, '', 1695285136, NULL);
INSERT INTO `dict_data` VALUES (87, 'sub_type', '刀片', '刀片', NULL, '', 'root', 1695285136, '', 1695285136, NULL);
INSERT INTO `dict_data` VALUES (88, 'outline_structure', '机架式', '机架式', NULL, '', 'root', 1695285215, '', 1695285215, NULL);
INSERT INTO `dict_data` VALUES (89, 'outline_structure', '刀箱', '刀箱', NULL, '', 'root', 1695285215, '', 1695285215, NULL);
INSERT INTO `dict_data` VALUES (90, 'outline_structure', '刀片', '刀片', NULL, '', 'root', 1695285215, '', 1695285215, NULL);
INSERT INTO `dict_data` VALUES (91, 'out_band_version', 'AMM', 'AMM', NULL, '', 'root', 1695285796, '', 1695285796, NULL);
INSERT INTO `dict_data` VALUES (92, 'producer-type', 'producer', '厂商', NULL, '', 'root', 1695286672, '', 1695286672, NULL);
INSERT INTO `dict_data` VALUES (93, 'producer-type', 'third_party_maintenance', '第三方维保服务商', NULL, '', 'root', 1695286672, '', 1695286672, NULL);
INSERT INTO `dict_data` VALUES (94, 'producer-type', 'supplier', '供应商', NULL, '', 'root', 1695286672, '', 1695286672, NULL);
INSERT INTO `dict_data` VALUES (95, 'producer-type', 'component_brand', '部件品牌', NULL, '', 'root', 1695286672, '', 1695286672, NULL);
INSERT INTO `dict_data` VALUES (96, 'basic_expansion', 'xingming-zhangsanasade', '姓名-张三阿萨德', 0, '', 'root', 1695366276, 'root', 1695723371, NULL);
INSERT INTO `dict_data` VALUES (97, 'basic_expansion', 'xingbie-nv', '性别-女', 5, '', 'root', 1695366276, 'root', 1695723384, NULL);
INSERT INTO `dict_data` VALUES (98, 'asset_expansion_fields', 'asset_expand_field', '资产扩展字段', NULL, '', 'root', 1695605034, '', 1695605034, NULL);
INSERT INTO `dict_data` VALUES (99, 'asset_expansion_fields', 'business_expand_field', '业务扩展字段', NULL, '', 'root', 1695605034, '', 1695605034, NULL);
INSERT INTO `dict_data` VALUES (100, 'asset_expansion_fields', 'common_asset_expand_field', '普通资产扩展字段', NULL, '', 'root', 1695605034, '', 1695605034, NULL);
INSERT INTO `dict_data` VALUES (102, 'link_method', 'telnet', 'telnet', 0, '', 'root', 1695696648, '', 1695696648, NULL);
INSERT INTO `dict_data` VALUES (103, 'link_method', 'ssh', 'ssh', 0, '', 'root', 1695696648, '', 1695696648, NULL);
INSERT INTO `dict_data` VALUES (108, 'asset_filter_status', 'alter-changed', '发生过变更的', 0, '', 'root', 1695720399, '', 1695720399, NULL);
INSERT INTO `dict_data` VALUES (109, 'asset_filter_status', 'maint-expired', '已过保修期的', 0, '', 'root', 1695720399, '', 1695720399, NULL);
INSERT INTO `dict_data` VALUES (110, 'asset_filter_status', 'status-1', '待上线', 0, '', 'root', 1695720399, '', 1695720399, NULL);
INSERT INTO `dict_data` VALUES (111, 'asset_filter_status', 'status-2', '已上线', 0, '', 'root', 1695720399, '', 1695720399, NULL);
INSERT INTO `dict_data` VALUES (117, 'maintenance_type', 'producer', '原厂维保', 0, '', 'root', 1695817242, '', 1695817242, NULL);
INSERT INTO `dict_data` VALUES (118, 'maintenance_type', 'third_party_maintenance', '第三方维保', 0, '', 'root', 1695817242, '', 1695817242, NULL);
INSERT INTO `dict_data` VALUES (119, 'offline_release_status', 'chear_business', '解除关联业务、解除检测、告警', 0, '', 'root', 1695883595, '', 1695883595, NULL);
INSERT INTO `dict_data` VALUES (120, 'offline_release_status', 'clear_production_ip', '清除生产IP', 0, '', 'root', 1695883595, '', 1695883595, NULL);
INSERT INTO `dict_data` VALUES (121, 'offline_release_status', 'clear_cabinet', '清除机柜位置', 0, '', 'root', 1695883595, '', 1695883595, NULL);
INSERT INTO `dict_data` VALUES (122, 'offline_release_status', 'clear_room', '清除所在机房', 0, '', 'root', 1695883595, '', 1695883595, NULL);
INSERT INTO `dict_data` VALUES (123, 'offline_release_status', 'clear_manage_ip', '清除管理IP', 0, '', 'root', 1695883595, '', 1695883595, NULL);
INSERT INTO `dict_data` VALUES (124, 'scrap_query_filter', 'all_scrap_device', '所有报废设备', 0, '', 'root', 1696750267, '', 1696750267, NULL);
INSERT INTO `dict_data` VALUES (125, 'scrap_query_filter', 'in_3_months_scrap', '近三个月内报废设备', 0, '', 'root', 1696750267, '', 1696750267, NULL);
INSERT INTO `dict_data` VALUES (126, 'scrap_query_filter', 'in_6_months_scrap', '近半年内报废设备', 0, '', 'root', 1696750267, '', 1696750267, NULL);
INSERT INTO `dict_data` VALUES (127, 'scrap_query_filter', 'in_this_year_scrap', '今年内报废设备', 0, '', 'root', 1696750267, '', 1696750267, NULL);
INSERT INTO `dict_data` VALUES (128, 'scrap_query_filter', 'in_1_year_scrap', '最近一年内报废设备', 0, '', 'root', 1696750267, '', 1696750267, NULL);

SET FOREIGN_KEY_CHECKS = 1;
