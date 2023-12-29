set names utf8mb4;

create database n9e_v6;
use n9e_v6;

CREATE TABLE `users` (
    `id` bigint unsigned not null auto_increment,
    `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'login name, cannot rename',
    `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'display name, chinese name',
    `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `phone` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `portrait` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'portrait image url',
    `roles` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'Admin | Standard | Guest, split by space',
    `status` int NOT NULL COMMENT '用户状态',
    `organization_id` int DEFAULT NULL COMMENT '组织id',
    `contacts` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci comment 'json e.g. {wecom:xx, dingtalk_robot_token:yy}',
    `maintainer` tinyint(1) not null default 0,
    `board_id` bigint unsigned not null default 0 COMMENT '默认首页模板id',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`username`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

INSERT INTO `users`(`id`, `username`, `nickname`, `password`, `phone`, `email`, `portrait`, `roles`, `status`, `organization_id`, `contacts`, `maintainer`, `create_at`, `create_by`, `update_at`, `update_by`, `board_id`) VALUES (1, 'admin', '超管', '042c05fffc2f49ca29a76223f3a41e83', '', '', '/image/avatar8.png', 'Admin', 1, 1, NULL, 0, 1698905269, 'system', 1703137997, 'admin', 1);

CREATE TABLE `user_group` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`create_by`),
    KEY (`update_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into user_group(id, name, create_at, create_by, update_at, update_by) values(1, '默认用户组', unix_timestamp(now()), 'admin', unix_timestamp(now()), 'admin');

CREATE TABLE `user_group_member` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint unsigned not null,
    `user_id` bigint unsigned not null,
    KEY (`group_id`),
    KEY (`user_id`),
    PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into user_group_member(group_id, user_id) values(1, 1);

CREATE TABLE `configs` (
    `id` bigint unsigned not null auto_increment,
    `ckey` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `cval` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`ckey`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `role` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into `role`(name, note) values('Admin', 'Administrator role');
insert into `role`(name, note) values('Standard', 'Ordinary user role');
insert into `role`(name, note) values('Guest', 'Readonly user role');

CREATE TABLE `role_operation`(
    `id` bigint unsigned not null auto_increment,
    `role_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `operation` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    KEY (`role_name`),
    KEY (`operation`),
    PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- Admin is special, who has no concrete operation but can do anything.
insert into `role_operation`(role_name, operation) values('Guest', '/metric/explorer');
insert into `role_operation`(role_name, operation) values('Guest', '/object/explorer');
insert into `role_operation`(role_name, operation) values('Guest', '/log/explorer');
insert into `role_operation`(role_name, operation) values('Guest', '/trace/explorer');
insert into `role_operation`(role_name, operation) values('Guest', '/help/version');
insert into `role_operation`(role_name, operation) values('Guest', '/help/contact');

insert into `role_operation`(role_name, operation) values('Standard', '/metric/explorer');
insert into `role_operation`(role_name, operation) values('Standard', '/object/explorer');
insert into `role_operation`(role_name, operation) values('Standard', '/log/explorer');
insert into `role_operation`(role_name, operation) values('Standard', '/trace/explorer');
insert into `role_operation`(role_name, operation) values('Standard', '/help/version');
insert into `role_operation`(role_name, operation) values('Standard', '/help/contact');
insert into `role_operation`(role_name, operation) values('Standard', '/help/servers');
insert into `role_operation`(role_name, operation) values('Standard', '/help/migrate');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-rules-built-in');
insert into `role_operation`(role_name, operation) values('Standard', '/dashboards-built-in');
insert into `role_operation`(role_name, operation) values('Standard', '/trace/dependencies');

insert into `role_operation`(role_name, operation) values('Admin', '/help/source');
insert into `role_operation`(role_name, operation) values('Admin', '/help/sso');
insert into `role_operation`(role_name, operation) values('Admin', '/help/notification-tpls');
insert into `role_operation`(role_name, operation) values('Admin', '/help/notification-settings');

insert into `role_operation`(role_name, operation) values('Standard', '/users');
insert into `role_operation`(role_name, operation) values('Standard', '/user-groups');
insert into `role_operation`(role_name, operation) values('Standard', '/user-groups/add');
insert into `role_operation`(role_name, operation) values('Standard', '/user-groups/put');
insert into `role_operation`(role_name, operation) values('Standard', '/user-groups/del');
insert into `role_operation`(role_name, operation) values('Standard', '/busi-groups');
insert into `role_operation`(role_name, operation) values('Standard', '/busi-groups/add');
insert into `role_operation`(role_name, operation) values('Standard', '/busi-groups/put');
insert into `role_operation`(role_name, operation) values('Standard', '/busi-groups/del');
insert into `role_operation`(role_name, operation) values('Standard', '/targets');
insert into `role_operation`(role_name, operation) values('Standard', '/targets/add');
insert into `role_operation`(role_name, operation) values('Standard', '/targets/put');
insert into `role_operation`(role_name, operation) values('Standard', '/targets/del');
insert into `role_operation`(role_name, operation) values('Standard', '/dashboards');
insert into `role_operation`(role_name, operation) values('Standard', '/dashboards/add');
insert into `role_operation`(role_name, operation) values('Standard', '/dashboards/put');
insert into `role_operation`(role_name, operation) values('Standard', '/dashboards/del');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-rules');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-rules/add');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-rules/put');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-rules/del');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-mutes');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-mutes/add');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-mutes/del');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-subscribes');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-subscribes/add');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-subscribes/put');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-subscribes/del');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-cur-events');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-cur-events/del');
insert into `role_operation`(role_name, operation) values('Standard', '/alert-his-events');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tpls');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tpls/add');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tpls/put');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tpls/del');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tasks');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tasks/add');
insert into `role_operation`(role_name, operation) values('Standard', '/job-tasks/put');
insert into `role_operation`(role_name, operation) values('Standard', '/recording-rules');
insert into `role_operation`(role_name, operation) values('Standard', '/recording-rules/add');
insert into `role_operation`(role_name, operation) values('Standard', '/recording-rules/put');
insert into `role_operation`(role_name, operation) values('Standard', '/recording-rules/del');

-- for alert_rule | collect_rule | mute | dashboard grouping
CREATE TABLE `busi_group` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `label_enable` tinyint(1) not null default 0,
    `label_value` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'if label_enable: label_value can not be blank',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into busi_group(id, name, create_at, create_by, update_at, update_by) values(1, '默认业务组', unix_timestamp(now()), 'root', unix_timestamp(now()), 'root');

CREATE TABLE `busi_group_member` (
    `id` bigint unsigned not null auto_increment,
    `busi_group_id` bigint not null comment 'busi group id',
    `user_group_id` bigint not null comment 'user group id',
    `perm_flag` char(2) not null comment 'ro | rw',
    PRIMARY KEY (`id`),
    KEY (`busi_group_id`),
    KEY (`user_group_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into busi_group_member(busi_group_id, user_group_id, perm_flag) values(1, 1, 'rw');

-- for dashboard new version
CREATE TABLE `board` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default 0 comment 'busi group id',
    `asset_id` bigint DEFAULT NULL,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `ident` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'split by space',
    `public` tinyint(1) not null default 0 comment '0:false 1:true',
    `built_in` tinyint(1) not null default 0 comment '0:false 1:true',
    `hide` tinyint(1) not null default 0 comment '0:false 1:true',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`group_id`, `name`),
    KEY(`ident`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of board
-- ----------------------------
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (1, 1, '首页', '', '', 0, 0, 0, 1700618966, 'root', 1703561005, 'admin', NULL);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (2, 1, 'Host Server', '', 'template', 0, 0, 0, 1702193566, 'admin', 1702256304, 'admin', 0);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (3, 1, 'MySQL', '', 'template', 0, 0, 0, 1702193713, 'admin', 1702195203, 'admin', 0);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (4, 1, 'Redis', '', 'template', 0, 0, 0, 1702199208, 'admin', 1702199844, 'admin', 0);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (5, 1, 'WebAPI', '', 'template', 0, 0, 0, 1702200143, 'admin', 1702204376, 'admin', 0);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (6, 1, 'TCP端口检测', '', 'template', 0, 0, 0, 1702204356, 'admin', 1702204604, 'admin', 0);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (7, 1, '网络端点', '', 'template', 0, 0, 0, 1702206481, 'admin', 1702206901, 'admin', 0);
INSERT INTO `board`(`id`, `group_id`, `name`, `ident`, `tags`, `public`, `built_in`, `hide`, `create_at`, `create_by`, `update_at`, `update_by`, `asset_id`) VALUES (8, 1, 'Switch', '', 'template', 0, 0, 0, 1702209071, 'admin', 1702210489, 'admin', 0);


-- for dashboard new version
CREATE TABLE `board_payload` (
    `id` bigint unsigned not null comment 'dashboard id',
    `payload` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `asset_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    UNIQUE KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of board_payload
-- ----------------------------
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (1, '{\"var\":[],\"panels\":[{\"type\":\"pie\",\"id\":\"c8a6ad98-183c-45d8-886c-71fb1a6ccf8b\",\"layout\":{\"h\":5,\"w\":8,\"x\":0,\"y\":0,\"i\":\"c8a6ad98-183c-45d8-886c-71fb1a6ccf8b\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"api\",\"datasourceValue\":999,\"targets\":[{\"refId\":\"A\",\"expr\":\"/api/n9e/api-service/2/execute\",\"legend\":\"{{name}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"资产健康度\",\"maxPerRow\":4,\"custom\":{\"angleField\":\"value\",\"colorField\":\"name\",\"calc\":\"lastNotNull\",\"legengPosition\":\"right\",\"donut\":true,\"labelWithName\":true,\"labelWithValue\":false,\"detailName\":\"详情\"},\"options\":{\"standardOptions\":{}},\"hidden\":false},{\"type\":\"pie\",\"id\":\"49cc582a-d756-48cd-8aa6-eca71c7b6900\",\"layout\":{\"h\":5,\"w\":8,\"x\":8,\"y\":0,\"i\":\"49c17483-4007-49ff-85e2-24a48d181a30\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"api\",\"datasourceValue\":999,\"targets\":[{\"refId\":\"A\",\"expr\":\"/api/n9e/api-service/3/execute\",\"legend\":\"{{name}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"资产监控状态\",\"maxPerRow\":4,\"custom\":{\"angleField\":\"value\",\"colorField\":\"name\",\"calc\":\"lastNotNull\",\"legengPosition\":\"right\",\"donut\":true,\"labelWithName\":true,\"labelWithValue\":false,\"detailName\":\"详情\"},\"options\":{\"standardOptions\":{\"decimals\":null}},\"hidden\":false},{\"type\":\"pie\",\"id\":\"21f88568-7ee4-4af0-ae4a-065768690ac4\",\"layout\":{\"h\":5,\"w\":8,\"x\":16,\"y\":0,\"i\":\"21f88568-7ee4-4af0-ae4a-065768690ac4\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":1,\"targets\":[{\"refId\":\"A\",\"expr\":\"topk(10, cpu_usage_active)and{cpu=\\\"cpu-total\\\"}\",\"legend\":\"{{agent_ip}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"CPU利用率Top10\",\"maxPerRow\":4,\"custom\":{\"calc\":\"lastNotNull\",\"legengPosition\":\"right\",\"detailName\":\"详情\"},\"options\":{\"standardOptions\":{}},\"hidden\":false},{\"type\":\"line\",\"id\":\"4292ba06-62c5-482c-8855-7812172a6b2c\",\"layout\":{\"h\":5,\"w\":8,\"x\":0,\"y\":5,\"i\":\"4292ba06-62c5-482c-8855-7812172a6b2c\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"api\",\"datasourceValue\":999,\"targets\":[{\"refId\":\"A\",\"expr\":\"/api/n9e/api-service/5/execute\",\"legend\":\"{{name}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"近7日告警趋势\",\"maxPerRow\":4,\"hidden\":false},{\"type\":\"column\",\"id\":\"6036e42b-d73a-4a57-938d-f32d59d83f0e\",\"layout\":{\"h\":5,\"w\":8,\"x\":8,\"y\":5,\"i\":\"6036e42b-d73a-4a57-938d-f32d59d83f0e\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"api\",\"datasourceValue\":999,\"targets\":[{\"refId\":\"A\",\"expr\":\"/api/n9e/api-service/4/execute\",\"legend\":\"{{name}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"资产分布状态\",\"maxPerRow\":4,\"custom\":{\"calc\":\"lastNotNull\",\"stack\":\"noraml\",\"seriesField\":\"type\"},\"options\":{\"standardOptions\":{\"util\":\"none\"}},\"hidden\":false},{\"type\":\"column\",\"id\":\"1b101282-95f6-442d-bbc4-2e838aecb5f0\",\"layout\":{\"h\":5,\"w\":8,\"x\":16,\"y\":5,\"i\":\"1b101282-95f6-442d-bbc4-2e838aecb5f0\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":1,\"targets\":[{\"refId\":\"A\",\"expr\":\"topk(10,mem_active)\",\"legend\":\"{{agent_ip}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"内存利用率Top10\",\"links\":[],\"maxPerRow\":4,\"custom\":{\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{}},\"hidden\":false},{\"type\":\"barGauge\",\"id\":\"c1744489-cf85-43ec-924c-9f20368d53bb\",\"layout\":{\"h\":5,\"w\":8,\"x\":0,\"y\":10,\"i\":\"c1744489-cf85-43ec-924c-9f20368d53bb\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":1,\"targets\":[{\"refId\":\"A\",\"expr\":\"rate(net_bits_recv[5m])\",\"legend\":\"{{instance}}-{{interface}}\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"网络监测趋势\",\"maxPerRow\":4,\"custom\":{\"calc\":\"lastNotNull\",\"baseColor\":\"#6395fa\",\"serieWidth\":20,\"sortOrder\":\"desc\"},\"options\":{\"standardOptions\":{}},\"hidden\":false},{\"type\":\"column\",\"id\":\"d941cae8-4c22-4434-afed-7bac0c37c58a\",\"layout\":{\"h\":5,\"w\":16,\"x\":8,\"y\":10,\"i\":\"d941cae8-4c22-4434-afed-7bac0c37c58a\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"api\",\"datasourceValue\":999,\"targets\":[{\"refId\":\"A\",\"expr\":\"/api/n9e/api-service/6/execute\",\"legend\":\"{{name}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"品牌故障率\",\"maxPerRow\":4,\"custom\":{\"calc\":\"lastNotNull\",\"stack\":\"off\"},\"options\":{\"standardOptions\":{\"util\":\"percentUnit\"}},\"hidden\":false},{\"type\":\"topo\",\"id\":\"32e33f5f-9d8c-4d10-96ac-4c22e510b2c0\",\"layout\":{\"h\":5,\"w\":16,\"x\":8,\"y\":15,\"i\":\"32e33f5f-9d8c-4d10-96ac-4c22e510b2c0\",\"isResizable\":true},\"version\":\"3.0.0\",\"custom\":{\"topo\":{\"cells\":[{\"position\":{\"x\":100,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"9959d5b6-4873-432e-94aa-4c3e381c05fc\",\"data\":{\"label\":\"服务器\"},\"zIndex\":1,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":187,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"e7eefc1e-7655-44e0-9c2f-bf067be45009\",\"data\":{\"label\":\"服务器\"},\"zIndex\":2,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":273,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"a42eec1a-bfbe-49eb-92b3-b33cb3eed012\",\"data\":{\"label\":\"服务器\"},\"zIndex\":3,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":365,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"f9ec4748-f769-4f9f-bca3-d043ca237bff\",\"data\":{\"label\":\"服务器\"},\"zIndex\":4,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":450,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"6bf099a1-5f52-4de8-8a72-885606ead602\",\"data\":{\"label\":\"服务器\"},\"zIndex\":5,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":538,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"1236238b-953d-4a4f-9fd2-f0e2b54d34e4\",\"data\":{\"label\":\"服务器\"},\"zIndex\":6,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":632,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"1e4512b7-8570-4a7e-986f-c78aac125e0b\",\"data\":{\"label\":\"服务器\"},\"zIndex\":7,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":718,\"y\":150},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Server\",\"id\":\"50f5406e-949d-4065-86b1-6fb583b8eafd\",\"data\":{\"label\":\"服务器\"},\"zIndex\":8,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":365,\"y\":10},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Switch\",\"id\":\"59c9e389-ac0e-4af8-8be4-04223ece6fd2\",\"data\":{\"label\":\"交换机\"},\"zIndex\":9,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"position\":{\"x\":450,\"y\":10},\"size\":{\"width\":40,\"height\":40},\"view\":\"react-shape-view\",\"attrs\":{\"body\":{\"fill\":\"#6796f5\",\"stroke\":\"#6796f5\",\"stroke-width\":0}},\"shape\":\"Switch\",\"id\":\"428e0a9e-7468-47c4-9354-ddca5c14d141\",\"data\":{\"label\":\"交换机\"},\"zIndex\":10,\"ports\":{\"groups\":{\"default\":{\"position\":\"center\",\"attrs\":{\"circle\":{\"magnet\":true,\"r\":5}}}},\"items\":[]}},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"a8efffd9-b17e-4f15-9618-dfc68b0daa51\",\"zIndex\":11,\"source\":{\"cell\":\"59c9e389-ac0e-4af8-8be4-04223ece6fd2\"},\"target\":{\"cell\":\"f9ec4748-f769-4f9f-bca3-d043ca237bff\"}},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"4a0121b5-5a71-4318-bacc-9075caee2c0c\",\"zIndex\":13,\"source\":{\"cell\":\"59c9e389-ac0e-4af8-8be4-04223ece6fd2\"},\"target\":{\"cell\":\"e7eefc1e-7655-44e0-9c2f-bf067be45009\"},\"vertices\":[{\"x\":385,\"y\":90},{\"x\":207,\"y\":90}]},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"ec648619-661b-4640-a58c-001329592ad5\",\"zIndex\":14,\"source\":{\"cell\":\"59c9e389-ac0e-4af8-8be4-04223ece6fd2\"},\"target\":{\"cell\":\"9959d5b6-4873-432e-94aa-4c3e381c05fc\"},\"vertices\":[{\"x\":385,\"y\":90},{\"x\":120,\"y\":90}]},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"eeb132a5-a0bc-45be-98f0-9d901936a58e\",\"zIndex\":15,\"source\":{\"cell\":\"59c9e389-ac0e-4af8-8be4-04223ece6fd2\"},\"target\":{\"cell\":\"a42eec1a-bfbe-49eb-92b3-b33cb3eed012\"},\"vertices\":[{\"x\":385,\"y\":90},{\"x\":293,\"y\":90}]},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"9a56fc3b-e0a6-4381-a0d0-e9d2541a72f2\",\"zIndex\":16,\"source\":{\"cell\":\"428e0a9e-7468-47c4-9354-ddca5c14d141\"},\"target\":{\"cell\":\"6bf099a1-5f52-4de8-8a72-885606ead602\"}},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"d672010b-b79a-4639-b1e0-a3197c357c92\",\"zIndex\":17,\"source\":{\"cell\":\"428e0a9e-7468-47c4-9354-ddca5c14d141\"},\"target\":{\"cell\":\"1236238b-953d-4a4f-9fd2-f0e2b54d34e4\"},\"vertices\":[{\"x\":470,\"y\":90},{\"x\":558,\"y\":90}]},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"79854ada-42ae-4bf5-a597-d498fc04d1ff\",\"zIndex\":18,\"source\":{\"cell\":\"428e0a9e-7468-47c4-9354-ddca5c14d141\"},\"target\":{\"cell\":\"1e4512b7-8570-4a7e-986f-c78aac125e0b\"},\"vertices\":[{\"x\":470,\"y\":90},{\"x\":652,\"y\":90}]},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"ddd13d66-d876-402b-8cb7-cd2155ff6d9e\",\"zIndex\":19,\"source\":{\"cell\":\"428e0a9e-7468-47c4-9354-ddca5c14d141\"},\"target\":{\"cell\":\"50f5406e-949d-4065-86b1-6fb583b8eafd\"},\"vertices\":[{\"x\":470,\"y\":90},{\"x\":738,\"y\":90}]},{\"shape\":\"edge\",\"attrs\":{\"line\":{\"stroke\":\"#00ba88\",\"targetMarker\":null}},\"id\":\"0d46d38f-6a22-46e4-8042-e5cc9aad8f43\",\"zIndex\":20,\"source\":{\"cell\":\"428e0a9e-7468-47c4-9354-ddca5c14d141\"},\"target\":{\"cell\":\"f9ec4748-f769-4f9f-bca3-d043ca237bff\"},\"vertices\":[{\"x\":470,\"y\":90},{\"x\":385,\"y\":90}]}]}},\"hidden\":false},{\"type\":\"row\",\"id\":\"e91023e6-5938-46da-8935-a08f278cfbc6\",\"name\":\"分组\",\"collapsed\":true,\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":20,\"i\":\"e91023e6-5938-46da-8935-a08f278cfbc6\",\"isResizable\":false},\"hidden\":false}],\"version\":\"3.0.0\",\"links\":[{\"title\":\"link\",\"url\":\"http://host.docker.internal:8888\",\"targetBlank\":false}]}', '');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (2, '{\"links\":[{\"targetBlank\":true,\"title\":\"n9e\",\"url\":\"https://n9e.github.io/\"},{\"targetBlank\":true,\"title\":\"author\",\"url\":\"http://flashcat.cloud/\"}],\"panels\":[{\"collapsed\":true,\"id\":\"2b2de3d1-65c8-4c39-9bea-02b754e0d751\",\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":0,\"i\":\"2b2de3d1-65c8-4c39-9bea-02b754e0d751\",\"isResizable\":false},\"name\":\"单机概况\",\"type\":\"row\",\"panels\":[]},{\"type\":\"stat\",\"id\":\"deec579b-3090-4344-a9a6-c1455c4a8e50\",\"layout\":{\"h\":4,\"w\":4,\"x\":0,\"y\":1,\"i\":\"deec579b-3090-4344-a9a6-c1455c4a8e50\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"system_uptime{asset_id=~\\\"$asset_id\\\"}/3600/24\",\"refId\":\"A\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"启动时长（单位：天）\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{\"value\":null}},\"options\":{\"valueMappings\":[],\"standardOptions\":{\"util\":\"none\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"gaugeN\",\"id\":\"7a7bd5db-d12e-49f0-92a8-15958e99ee54\",\"layout\":{\"h\":4,\"w\":4,\"x\":4,\"y\":1,\"i\":\"7a7bd5db-d12e-49f0-92a8-15958e99ee54\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"100-cpu_usage_idle{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"}\",\"refId\":\"A\",\"legend\":\"{{cpu}}\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"CPU使用率\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{\"util\":\"percent\"},\"thresholds\":{\"steps\":[{\"color\":\"#3FC453\",\"value\":null,\"type\":\"base\"},{\"color\":\"#FF9919\",\"value\":60},{\"color\":\"#FF656B\",\"value\":80}]}}},{\"type\":\"gaugeN\",\"id\":\"8a814265-54ad-419c-8cb7-e1f84a242de0\",\"layout\":{\"h\":4,\"w\":4,\"x\":8,\"y\":1,\"i\":\"8a814265-54ad-419c-8cb7-e1f84a242de0\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"mem_used_percent{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"A\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"内存使用率\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{\"util\":\"percent\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#3FC453\",\"value\":null,\"type\":\"base\"},{\"color\":\"#FF9919\",\"value\":60},{\"color\":\"#FF656B\",\"value\":80}]}}},{\"type\":\"gaugeN\",\"id\":\"ebef1143-4203-4843-b8fd-67a172d344a6\",\"layout\":{\"h\":4,\"w\":4,\"x\":12,\"y\":1,\"i\":\"ebef1143-4203-4843-b8fd-67a172d344a6\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"topk(1, disk_used_percent{asset_id=\\\"${asset_id}\\\"})\",\"legend\":\"{{path}}\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"最大分区使用率\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"valueAndName\",\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{\"util\":\"percent\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#3FC453\",\"value\":null,\"type\":\"base\"},{\"color\":\"#FF9919\",\"value\":70},{\"color\":\"#FF656B\",\"value\":85}]}}},{\"type\":\"stat\",\"id\":\"d7d11972-5c5b-4bc6-98f8-bbbe9f018896\",\"layout\":{\"h\":4,\"w\":4,\"x\":16,\"y\":1,\"i\":\"d7d11972-5c5b-4bc6-98f8-bbbe9f018896\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(net_bits_recv{asset_id=\\\"${asset_id}\\\"}[1m]))\",\"refId\":\"A\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"入流量(b/s)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"area\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{\"value\":null}},\"options\":{\"valueMappings\":[{\"match\":{\"from\":0,\"to\":50},\"result\":{\"color\":\"#129b22\"},\"type\":\"range\"},{\"match\":{\"from\":50,\"to\":100},\"result\":{\"color\":\"#f51919\"},\"type\":\"range\"}],\"standardOptions\":{\"util\":\"bitsIEC\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"8e919f07-c87f-431d-a0bd-6005402f2a99\",\"layout\":{\"h\":4,\"w\":4,\"x\":20,\"y\":1,\"i\":\"81190f39-889b-4d32-9de0-c9b99a4092ea\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(net_bits_sent{asset_id=\\\"${asset_id}\\\"}[1m]))\",\"refId\":\"A\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"出流量(b/s)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"area\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{\"value\":null}},\"options\":{\"valueMappings\":[{\"match\":{\"from\":0,\"to\":50},\"result\":{\"color\":\"#129b22\"},\"type\":\"range\"},{\"match\":{\"from\":50,\"to\":100},\"result\":{\"color\":\"#f51919\"},\"type\":\"range\"}],\"standardOptions\":{\"util\":\"bitsIEC\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"collapsed\":true,\"id\":\"aabb8263-1a9b-43fb-bee1-6c532f5012a3\",\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":5,\"i\":\"aabb8263-1a9b-43fb-bee1-6c532f5012a3\",\"isResizable\":false},\"name\":\"CPU\",\"type\":\"row\",\"panels\":[]},{\"type\":\"timeseriesN\",\"id\":\"1559d880-7e26-4e42-9427-4e55fb6f67be\",\"layout\":{\"h\":5,\"w\":8,\"x\":0,\"y\":6,\"i\":\"1559d880-7e26-4e42-9427-4e55fb6f67be\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"avg by (cpu) (cpu_usage_active{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"legend\":\"{{cpu}}\",\"refId\":\"A\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"CPU使用率\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"list\",\"placement\":\"bottom\"},\"standardOptions\":{\"util\":\"percent\",\"min\":0,\"max\":null,\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"linear\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"043c26de-d19f-4fe8-a615-2b7c10ceb828\",\"layout\":{\"h\":5,\"w\":8,\"x\":8,\"y\":6,\"i\":\"043c26de-d19f-4fe8-a615-2b7c10ceb828\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"avg(cpu_usage_guest{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"legend\":\"guest\",\"refId\":\"A\",\"step\":null},{\"expr\":\"avg(cpu_usage_iowait{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"legend\":\"iowait\",\"refId\":\"B\",\"step\":null},{\"expr\":\"avg(cpu_usage_user{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"refId\":\"C\",\"legend\":\"user\",\"step\":null},{\"expr\":\"avg(cpu_usage_system{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"refId\":\"D\",\"legend\":\"system\",\"step\":null},{\"expr\":\"avg(cpu_usage_irq{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"refId\":\"E\",\"legend\":\"irq\"},{\"expr\":\"avg(cpu_usage_softirq{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"refId\":\"F\",\"legend\":\"softirq\",\"step\":null},{\"expr\":\"avg(cpu_usage_nice{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"refId\":\"G\",\"legend\":\"nice\",\"step\":null},{\"expr\":\"avg(cpu_usage_steal{asset_id=~\\\"$asset_id\\\",cpu=\\\"cpu-total\\\"})[1m]\",\"refId\":\"H\",\"legend\":\"steal\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"CPU使用率详情\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"percent\",\"min\":null,\"max\":null,\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"linear\",\"fillOpacity\":0.5,\"stack\":\"noraml\"}},{\"type\":\"timeseriesN\",\"id\":\"a420ce25-6968-47f8-8335-60cde70fd062\",\"layout\":{\"h\":5,\"w\":8,\"x\":16,\"y\":6,\"i\":\"a420ce25-6968-47f8-8335-60cde70fd062\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"system_load15{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"A\",\"legend\":\"system_load15\",\"step\":null},{\"expr\":\"system_load1{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"B\",\"legend\":\"system_load1\",\"step\":null},{\"expr\":\"system_load5{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"C\",\"legend\":\"system_load5\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"CPU负载\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"linear\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"b7a3c99f-a796-4b76-89b5-cbddd566f91c\",\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":11,\"i\":\"b7a3c99f-a796-4b76-89b5-cbddd566f91c\",\"isResizable\":false},\"name\":\"内存详情\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"239aacdf-1982-428b-b240-57f4ce7f946d\",\"layout\":{\"h\":5,\"w\":12,\"x\":0,\"y\":12,\"i\":\"239aacdf-1982-428b-b240-57f4ce7f946d\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"mem_used_percent{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"A\",\"step\":null,\"legend\":\"使用率\"},{\"expr\":\"mem_cached{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"B\",\"step\":null,\"legend\":\"cached\"},{\"expr\":\"mem_buffered{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"C\",\"step\":null,\"legend\":\"bufferd\"},{\"expr\":\"mem_inactive{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"D\",\"step\":null,\"legend\":\"inactive\"},{\"expr\":\"mem_mapped{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"E\",\"step\":null,\"legend\":\"mapped\"},{\"expr\":\"mem_shared{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"F\",\"step\":null,\"legend\":\"shared\"},{\"expr\":\"mem_swap_cached{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"G\",\"step\":null,\"legend\":\"swap_cached\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"用户态内存使用\",\"description\":\"内存指标可参考链接 [/PROC/MEMINFO之谜](http://linuxperf.com/?p=142) \",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"percent\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"00ed6e4d-c979-4938-a20e-56d42ca452cf\",\"layout\":{\"h\":5,\"w\":12,\"x\":12,\"y\":12,\"i\":\"00ed6e4d-c979-4938-a20e-56d42ca452cf\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"mem_slab{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"A\",\"legend\":\"slab\",\"step\":null},{\"expr\":\"mem_sreclaimable{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"B\",\"legend\":\"sreclaimable\",\"step\":null},{\"expr\":\"mem_sunreclaim{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"C\",\"legend\":\"sunreclaim\",\"step\":null},{\"expr\":\"mem_vmalloc_used{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"D\",\"legend\":\"vmalloc_used\",\"step\":null},{\"expr\":\"mem_vmalloc_chunk{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"E\",\"legend\":\"vmalloc_chunk\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"内核态内存使用\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bytesIEC\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"842a8c48-0e93-40bf-8f28-1b2f837e5c19\",\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":17,\"i\":\"842a8c48-0e93-40bf-8f28-1b2f837e5c19\",\"isResizable\":false},\"name\":\"磁盘详情\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"9f416607-5a65-4dd9-a5c8-30a61bbd7ae7\",\"layout\":{\"h\":5,\"w\":8,\"x\":0,\"y\":18,\"i\":\"9f416607-5a65-4dd9-a5c8-30a61bbd7ae7\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"disk_total{asset_id=\\\"${asset_id}\\\"}\",\"legend\":\"总量-{{path}}\"},{\"expr\":\"disk_used{asset_id=\\\"${asset_id}\\\"}\",\"refId\":\"B\",\"legend\":\"已用-{{path}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"磁盘使用\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bytesIEC\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"d825671f-7dc5-46a2-89dc-4fff084a3ae0\",\"layout\":{\"h\":5,\"w\":8,\"x\":8,\"y\":18,\"i\":\"d825671f-7dc5-46a2-89dc-4fff084a3ae0\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"linux_sysctl_fs_file_max{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"A\",\"legend\":\"file_max\",\"step\":null},{\"expr\":\"linux_sysctl_fs_file_nr{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"B\",\"legend\":\"file_nr\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"fd使用\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"d27b522f-9c70-42f2-9e31-fed3816fd675\",\"layout\":{\"h\":5,\"w\":8,\"x\":16,\"y\":18,\"i\":\"d27b522f-9c70-42f2-9e31-fed3816fd675\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"disk_inodes_total{asset_id=~\\\"$asset_id\\\",path!~\\\"/var.*\\\"}\",\"legend\":\"inodes_total_{{path}}\",\"refId\":\"A\",\"step\":null},{\"expr\":\"disk_inodes_used{asset_id=~\\\"$asset_id\\\",path!~\\\"/var.*\\\"}\",\"legend\":\"inodes_used_{{path}}\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"inode\",\"description\":\"windows不适用\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"bbd1ebda-99f6-419c-90a5-5f84973976dd\",\"layout\":{\"h\":5,\"w\":8,\"x\":0,\"y\":23,\"i\":\"bbd1ebda-99f6-419c-90a5-5f84973976dd\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(diskio_read_bytes{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{name}}-读\",\"refId\":\"A\",\"step\":null},{\"expr\":\"rate(diskio_write_bytes{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{name}}-写\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"IO吞吐量\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bytesIEC\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"f645741e-c632-4685-b267-c7ad26b5c10e\",\"layout\":{\"h\":5,\"w\":8,\"x\":8,\"y\":23,\"i\":\"f645741e-c632-4685-b267-c7ad26b5c10e\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(diskio_reads{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{name}}-读\",\"refId\":\"A\",\"step\":null},{\"expr\":\"rate(diskio_writes{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{name}}-写\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"IOPS\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"none\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"d6b45598-54c6-4b36-a896-0a7529ac21f8\",\"layout\":{\"h\":5,\"w\":8,\"x\":16,\"y\":23,\"i\":\"d6b45598-54c6-4b36-a896-0a7529ac21f8\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(diskio_write_time{asset_id=~\\\"$asset_id\\\"}[1m])/rate(diskio_writes{asset_id=~\\\"$asset_id\\\"}[1m])+rate(diskio_read_time{asset_id=~\\\"$asset_id\\\"}[1m])/rate(diskio_reads{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{name}}\",\"refId\":\"A\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"iowait\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"linear\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"307152d2-708c-4736-98cf-08b886cbf7f2\",\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":28,\"i\":\"307152d2-708c-4736-98cf-08b886cbf7f2\",\"isResizable\":false},\"name\":\"网络详情\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"f2ee5d32-737c-4095-b6b7-b15b778ffdb9\",\"layout\":{\"h\":5,\"w\":6,\"x\":0,\"y\":29,\"i\":\"f2ee5d32-737c-4095-b6b7-b15b778ffdb9\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(net_bits_recv{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-接收\",\"refId\":\"A\",\"step\":null},{\"expr\":\"rate(net_bits_sent{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-发送\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"网络流量\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bitsIEC\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"9113323a-98f5-4bff-a8ce-3b459e7e2190\",\"layout\":{\"h\":5,\"w\":6,\"x\":6,\"y\":29,\"i\":\"9113323a-98f5-4bff-a8ce-3b459e7e2190\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(net_packets_recv{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-接收\",\"refId\":\"A\",\"step\":null},{\"expr\":\"rate(net_packets_sent{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-发送\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"网络数据包\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"9634c41c-e124-4d7f-9406-0f86753e8d70\",\"layout\":{\"h\":5,\"w\":6,\"x\":12,\"y\":29,\"i\":\"9634c41c-e124-4d7f-9406-0f86753e8d70\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(net_err_in{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-入方向\",\"refId\":\"A\",\"step\":null},{\"expr\":\"rate(net_err_out{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-出方向\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"错包\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"4123f4c1-bf8e-400e-b267-8d7f6a92691a\",\"layout\":{\"h\":5,\"w\":6,\"x\":18,\"y\":29,\"i\":\"4123f4c1-bf8e-400e-b267-8d7f6a92691a\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(net_drop_in{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-入方向\",\"refId\":\"A\",\"step\":null},{\"expr\":\"rate(net_drop_out{asset_id=~\\\"$asset_id\\\"}[1m])\",\"legend\":\"{{interface}}-出方向\",\"refId\":\"B\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"丢包\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"cfb80689-de7b-47fb-9155-052b796dd7f5\",\"layout\":{\"h\":5,\"w\":24,\"x\":0,\"y\":34,\"i\":\"cfb80689-de7b-47fb-9155-052b796dd7f5\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"netstat_tcp_established{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"A\",\"legend\":\"established\",\"step\":null},{\"expr\":\"netstat_tcp_listen{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"B\",\"legend\":\"listen\",\"step\":null},{\"expr\":\"netstat_tcp_time_wait{asset_id=~\\\"$asset_id\\\"}\",\"refId\":\"C\",\"legend\":\"time_wait\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"tcp连接\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"row\",\"id\":\"613bbbb0-98e3-4e5d-9b5a-c3a54c231105\",\"name\":\"进程信息\",\"collapsed\":true,\"layout\":{\"h\":1,\"w\":24,\"x\":0,\"y\":39,\"i\":\"613bbbb0-98e3-4e5d-9b5a-c3a54c231105\",\"isResizable\":false}},{\"type\":\"table\",\"id\":\"5bdbd6c4-8aa6-4e12-97e4-9529fa1ca17c\",\"layout\":{\"h\":7,\"w\":24,\"x\":0,\"y\":40,\"i\":\"5bdbd6c4-8aa6-4e12-97e4-9529fa1ca17c\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"label_join(procstat_cpu_usage, \\\"pname\\\", \\\"-\\\", \\\"pid\\\", \\\"comm\\\")\",\"instant\":true,\"legend\":\"cpu使用率\"},{\"expr\":\"label_join(procstat_mem_usage, \\\"pname\\\", \\\"-\\\", \\\"pid\\\", \\\"comm\\\")\",\"refId\":\"B\",\"legend\":\"内存使用率\",\"instant\":true},{\"expr\":\"label_join(procstat_num_threads, \\\"pname\\\", \\\"-\\\", \\\"pid\\\", \\\"comm\\\")\",\"refId\":\"C\",\"legend\":\"线程数\",\"instant\":true},{\"expr\":\"label_join(rate(procstat_read_bytes[5m]), \\\"pname\\\", \\\"-\\\", \\\"pid\\\", \\\"comm\\\")\",\"refId\":\"E\",\"instant\":true,\"legend\":\"IO读\"},{\"expr\":\"label_join(rate(procstat_write_bytes[5m]), \\\"pname\\\", \\\"-\\\", \\\"pid\\\", \\\"comm\\\")\",\"refId\":\"F\",\"instant\":true,\"legend\":\"IO写\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{\"indexByName\":{\"pid\":0},\"excludeByName\":{\"pid\":false},\"renameByName\":{\"pid\":\"\",\"comm\":\"进程名称\"}}}],\"name\":\"进程列表\",\"maxPerRow\":4,\"custom\":{\"showHeader\":true,\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"displayMode\":\"labelValuesToRows\",\"aggrDimension\":\"pname\",\"sortColumn\":\"pname\",\"sortOrder\":\"ascend\"},\"options\":{\"standardOptions\":{\"util\":\"none\"}},\"overrides\":[{\"matcher\":{\"value\":\"A\"},\"properties\":{\"standardOptions\":{\"util\":\"percent\",\"decimals\":2}}},{\"type\":\"special\",\"matcher\":{\"value\":\"B\"},\"properties\":{\"standardOptions\":{\"util\":\"percent\",\"decimals\":2}}},{\"type\":\"special\",\"matcher\":{\"value\":\"E\"},\"properties\":{\"standardOptions\":{\"util\":\"bytesIEC\"}}},{\"type\":\"special\",\"matcher\":{\"value\":\"F\"},\"properties\":{\"standardOptions\":{\"util\":\"bytesIEC\"}}}]}],\"var\":[{\"name\":\"prom\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${prom}\"},\"definition\":\"label_values(asset_up, asset_id)\"}],\"version\":\"3.0.0\"}', '主机设备');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (3, '{\"panels\":[{\"collapsed\":true,\"id\":\"fe0e2a5d-4e82-4eaf-b13a-6d98aa6b6860\",\"layout\":{\"h\":1,\"i\":\"fe0e2a5d-4e82-4eaf-b13a-6d98aa6b6860\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":0},\"name\":\"基本信息\",\"type\":\"row\"},{\"type\":\"stat\",\"id\":\"80079949-dbff-48fe-a1eb-54b646c30135\",\"layout\":{\"h\":3,\"i\":\"80079949-dbff-48fe-a1eb-54b646c30135\",\"isResizable\":true,\"w\":6,\"x\":0,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"min(mysql_global_status_uptime{asset_id=\\\"$asset_id\\\"})/3600/24\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"启动时长(单位: 天)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"valueMappings\":[{\"type\":\"range\",\"result\":{\"color\":\"#ffae39\"},\"match\":{\"to\":1}}],\"standardOptions\":{\"util\":\"none\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"9fd6dd09-d131-4c0e-88ea-ed62c72baf97\",\"layout\":{\"h\":3,\"i\":\"9fd6dd09-d131-4c0e-88ea-ed62c72baf97\",\"isResizable\":true,\"w\":6,\"x\":6,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(mysql_global_status_queries{asset_id=\\\"$asset_id\\\"}[5m])\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"每秒查询次数(QPS)\",\"description\":\"mysql_global_status_queries\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"valueMappings\":[{\"type\":\"range\",\"result\":{\"color\":\"#ff9919\"},\"match\":{\"from\":100}}],\"standardOptions\":{\"decimals\":2},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"24913190-b86d-44b7-a8db-555351d9d3c2\",\"layout\":{\"h\":3,\"i\":\"24913190-b86d-44b7-a8db-555351d9d3c2\",\"isResizable\":true,\"w\":6,\"x\":12,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"avg(mysql_global_variables_innodb_buffer_pool_size{asset_id=\\\"$asset_id\\\"})\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"缓存池空间(InnoDB Buffer Pool)\",\"description\":\"\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{\"util\":\"bytesIEC\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"94a1e97e-2241-4e05-a9e9-a9b1e69d1070\",\"layout\":{\"h\":3,\"i\":\"94a1e97e-2241-4e05-a9e9-a9b1e69d1070\",\"isResizable\":true,\"w\":6,\"x\":18,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(increase(mysql_global_status_table_locks_waited{asset_id=\\\"$asset_id\\\"}[5m]))\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"表锁等待(5min)\",\"description\":\"**Table Locks**\\n\\nMySQL takes a number of different locks for varying reasons. In this graph we see how many Table level locks MySQL has requested from the storage engine. In the case of InnoDB, many times the locks could actually be row locks as it only takes table level locks in a few specific cases.\\n\\nIt is most useful to compare Locks Immediate and Locks Waited. If Locks waited is rising, it means you have lock contention. Otherwise, Locks Immediate rising and falling is normal activity.\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"valueMappings\":[{\"match\":{\"from\":1},\"result\":{\"color\":\"#e70d0d\"},\"type\":\"range\"}],\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"collapsed\":true,\"id\":\"ca82d30f-8e0d-4caa-8a00-2ed9efe4ad85\",\"layout\":{\"h\":1,\"i\":\"ca82d30f-8e0d-4caa-8a00-2ed9efe4ad85\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":4},\"name\":\"连接信息\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"e2c85e72-0286-49bc-8ddb-5fba5f449b53\",\"layout\":{\"h\":7,\"i\":\"e2c85e72-0286-49bc-8ddb-5fba5f449b53\",\"isResizable\":true,\"w\":12,\"x\":0,\"y\":5},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(mysql_global_status_threads_connected{asset_id=\\\"$asset_id\\\"})\",\"legend\":\"连接数(Connections)\"},{\"expr\":\"sum(mysql_global_status_max_used_connections{asset_id=\\\"$asset_id\\\"})\",\"legend\":\"最大使用连接数(Max Used Connections)\"},{\"expr\":\"sum(mysql_global_variables_max_connections{asset_id=\\\"$asset_id\\\"})\",\"legend\":\"最大连接数(Max Connections)\"},{\"expr\":\"sum(rate(mysql_global_status_aborted_connects{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"中断连接数(Aborted Connections)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"客户端连接(MySQL Connections)\",\"description\":\"**Max Connections** \\n\\nMax Connections is the maximum permitted number of simultaneous client connections. By default, this is 151. Increasing this value increases the number of file descriptors that mysqld requires. If the required number of descriptors are not available, the server reduces the value of Max Connections.\\n\\nmysqld actually permits Max Connections + 1 clients to connect. The extra connection is reserved for use by accounts that have the SUPER privilege, such as root.\\n\\nMax Used Connections is the maximum number of connections that have been in use simultaneously since the server started.\\n\\nConnections is the number of connection attempts (successful or not) to the MySQL server.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"fbd43ac2-159d-4e55-8bc6-800d1bbfbd59\",\"layout\":{\"h\":7,\"i\":\"fbd43ac2-159d-4e55-8bc6-800d1bbfbd59\",\"isResizable\":true,\"w\":12,\"x\":12,\"y\":5},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(mysql_global_status_threads_connected{asset_id=\\\"$asset_id\\\"})\",\"legend\":\"已连接(Threads Connected)\"},{\"expr\":\"sum(mysql_global_status_threads_running{asset_id=\\\"$asset_id\\\"})\",\"legend\":\"运行中(Threads Running)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"客户端线程(MySQL Client Thread Activity)\",\"description\":\"Threads Connected is the number of open connections, while Threads Running is the number of threads not sleeping.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"cb81def4-ac63-4d42-b66e-440f9061794b\",\"layout\":{\"h\":1,\"i\":\"cb81def4-ac63-4d42-b66e-440f9061794b\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":12},\"name\":\"查询性能\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"5fa65a30-a49b-457f-b46a-11d2029188bd\",\"layout\":{\"h\":7,\"i\":\"5fa65a30-a49b-457f-b46a-11d2029188bd\",\"isResizable\":true,\"w\":12,\"x\":0,\"y\":13},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(mysql_global_status_created_tmp_tables{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"临时表(Created Tmp Tables)\"},{\"expr\":\"sum(rate(mysql_global_status_created_tmp_disk_tables{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"临时磁盘表(Created Tmp Disk Tables)\"},{\"expr\":\"sum(rate(mysql_global_status_created_tmp_files{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"临时文件(Created Tmp Files)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"临时对象(MySQL Temporary Objects)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"20efd251-6207-4cec-aa3b-4351e8e9b125\",\"layout\":{\"h\":7,\"i\":\"20efd251-6207-4cec-aa3b-4351e8e9b125\",\"isResizable\":true,\"w\":12,\"x\":12,\"y\":13},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(mysql_global_status_select_full_join{ asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Select Full Join\"},{\"expr\":\"sum(rate(mysql_global_status_select_full_range_join{ asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Select Full Range Join\"},{\"expr\":\"sum(rate(mysql_global_status_select_range{ asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Select Range\"},{\"expr\":\"sum(rate(mysql_global_status_select_range_check{ asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Select Range Check\"},{\"expr\":\"sum(rate(mysql_global_status_select_scan{ asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Select Scan\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"查询类型(MySQL Select Types)\",\"description\":\"**MySQL Select Types**\\n\\nAs with most relational databases, selecting based on indexes is more efficient than scanning an entire table\'s data. Here we see the counters for selects not done with indexes.\\n\\n* ***Select Scan*** is how many queries caused full table scans, in which all the data in the table had to be read and either discarded or returned.\\n* ***Select Range*** is how many queries used a range scan, which means MySQL scanned all rows in a given range.\\n* ***Select Full Join*** is the number of joins that are not joined on an index, this is usually a huge performance hit.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"a4d0c5fb-04e0-4627-8722-ae996d70e2aa\",\"layout\":{\"h\":7,\"i\":\"a4d0c5fb-04e0-4627-8722-ae996d70e2aa\",\"isResizable\":true,\"w\":12,\"x\":0,\"y\":20},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(mysql_global_status_sort_rows{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Sort Rows\"},{\"expr\":\"sum(rate(mysql_global_status_sort_range{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Sort Range\"},{\"expr\":\"sum(rate(mysql_global_status_sort_merge_passes{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Sort Merge Passes\"},{\"expr\":\"sum(rate(mysql_global_status_sort_scan{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Sort Scan\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"排序(MySQL Sorts)\",\"description\":\"**MySQL Sorts**\\n\\nDue to a query\'s structure, order, or other requirements, MySQL sorts the rows before returning them. For example, if a table is ordered 1 to 10 but you want the results reversed, MySQL then has to sort the rows to return 10 to 1.\\n\\nThis graph also shows when sorts had to scan a whole table or a given range of a table in order to return the results and which could not have been sorted via an index.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"2e13ada4-1128-440d-9360-028f16c3779b\",\"layout\":{\"h\":7,\"i\":\"2e13ada4-1128-440d-9360-028f16c3779b\",\"isResizable\":true,\"w\":12,\"x\":12,\"y\":20},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(mysql_global_status_slow_queries{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"Slow Queries\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"慢查询(MySQL Slow Queries)\",\"description\":\"**MySQL Slow Queries**\\n\\nSlow queries are defined as queries being slower than the long_query_time setting. For example, if you have long_query_time set to 3, all queries that take longer than 3 seconds to complete will show on this graph.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"c9df805c-8ae7-41d7-b28b-575f478fd9ce\",\"layout\":{\"h\":1,\"i\":\"c9df805c-8ae7-41d7-b28b-575f478fd9ce\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":27},\"name\":\"网络\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"6107714f-bedd-437c-b6e4-d6eb74db6d30\",\"layout\":{\"h\":7,\"i\":\"6107714f-bedd-437c-b6e4-d6eb74db6d30\",\"isResizable\":true,\"w\":24,\"x\":0,\"y\":28},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(mysql_global_status_bytes_received{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"入流量(Inbound)\"},{\"expr\":\"sum(rate(mysql_global_status_bytes_sent{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"出流量(Outbound)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"网络流量(MySQL Network Traffic)\",\"description\":\"**MySQL Network Traffic**\\n\\nHere we can see how much network traffic is generated by MySQL. Outbound is network traffic sent from MySQL and Inbound is network traffic MySQL has received.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"00fd2b70-a133-4ad7-bd56-69a3c91ecf0c\",\"layout\":{\"h\":1,\"i\":\"00fd2b70-a133-4ad7-bd56-69a3c91ecf0c\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":35},\"name\":\"命令处理\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"f90ca2bc-0809-45f6-88b6-e258805def04\",\"layout\":{\"h\":7,\"i\":\"f90ca2bc-0809-45f6-88b6-e258805def04\",\"isResizable\":true,\"w\":24,\"x\":0,\"y\":36},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"topk(10, rate(mysql_global_status_commands_total{asset_id=\\\"$asset_id\\\"}[5m])>0)\",\"legend\":\"Com_{{command}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"命令执行(Top Command Counters)\",\"description\":\"**Top Command Counters**\\n\\nThe Com_{{xxx}} statement counter variables indicate the number of times each xxx statement has been executed. There is one status variable for each type of statement. For example, Com_delete and Com_update count [``DELETE``](https://dev.mysql.com/doc/refman/5.7/en/delete.html) and [``UPDATE``](https://dev.mysql.com/doc/refman/5.7/en/update.html) statements, respectively. Com_delete_multi and Com_update_multi are similar but apply to [``DELETE``](https://dev.mysql.com/doc/refman/5.7/en/delete.html) and [``UPDATE``](https://dev.mysql.com/doc/refman/5.7/en/update.html) statements that use multiple-table syntax.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"74e1844d-a918-48fa-a29f-6535dc087dac\",\"layout\":{\"h\":7,\"i\":\"74e1844d-a918-48fa-a29f-6535dc087dac\",\"isResizable\":true,\"w\":12,\"x\":0,\"y\":43},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(mysql_global_status_handlers_total{asset_id=\\\"$asset_id\\\", handler!~\\\"commit|rollback|savepoint.*|prepare\\\"}[5m])\",\"legend\":\"{{handler}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"处理器(MySQL Handlers)\",\"description\":\"**MySQL Handlers**\\n\\nHandler statistics are internal statistics on how MySQL is selecting, updating, inserting, and modifying rows, tables, and indexes.\\n\\nThis is in fact the layer between the Storage Engine and MySQL.\\n\\n* `read_rnd_next` is incremented when the server performs a full table scan and this is a counter you don\'t really want to see with a high value.\\n* `read_key` is incremented when a read is done with an index.\\n* `read_next` is incremented when the storage engine is asked to \'read the next index entry\'. A high value means a lot of index scans are being done.\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"b2c3a13d-898f-407b-b6a9-db852072b12f\",\"layout\":{\"h\":7,\"i\":\"b2c3a13d-898f-407b-b6a9-db852072b12f\",\"isResizable\":true,\"w\":12,\"x\":12,\"y\":43},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(mysql_global_status_handlers_total{asset_id=\\\"$asset_id\\\", handler=~\\\"commit|rollback|savepoint.*|prepare\\\"}[5m])\",\"legend\":\"{{handler}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"事务(MySQL Transaction Handlers)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"c32a02da-6c61-4b9e-9365-c0b56088fabc\",\"layout\":{\"h\":1,\"i\":\"c32a02da-6c61-4b9e-9365-c0b56088fabc\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":50},\"name\":\"打开文件信息\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"fc13eadb-890d-4184-ac16-943d54188db8\",\"layout\":{\"h\":7,\"i\":\"fc13eadb-890d-4184-ac16-943d54188db8\",\"isResizable\":true,\"w\":24,\"x\":0,\"y\":51},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"mysql_global_variables_open_files_limit{asset_id=\\\"$asset_id\\\"}\",\"legend\":\"文件打开限制(Open Files Limit)\"},{\"expr\":\"mysql_global_status_open_files{asset_id=\\\"$asset_id\\\"}\",\"legend\":\"打开文件数量(Open Files)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"打开文件数量(MySQL Open Files)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"6f596e65-3e4b-4d9a-aad7-a32c8c7b8239\",\"layout\":{\"h\":1,\"i\":\"6f596e65-3e4b-4d9a-aad7-a32c8c7b8239\",\"isResizable\":false,\"w\":24,\"x\":0,\"y\":58},\"name\":\"打开表信息\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"0b78fbb5-a0b4-4a1b-98b1-af15bc91779d\",\"layout\":{\"h\":7,\"i\":\"0b78fbb5-a0b4-4a1b-98b1-af15bc91779d\",\"isResizable\":true,\"w\":12,\"x\":0,\"y\":59},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(mysql_global_status_table_open_cache_hits{asset_id=\\\"$asset_id\\\"}[5m])\\n/\\n(\\nrate(mysql_global_status_table_open_cache_hits{asset_id=\\\"$asset_id\\\"}[5m])\\n+\\nrate(mysql_global_status_table_open_cache_misses{asset_id=\\\"$asset_id\\\"}[5m])\\n)\",\"legend\":\"表缓存命中率(Table Open Cache Hit Ratio)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"表缓存命中率(Table Open Cache Hit Ratio  Mysql 5.6.6+)\",\"description\":\"**MySQL Table Open Cache Status**\\n\\nThe recommendation is to set the `table_open_cache_instances` to a loose correlation to virtual CPUs, keeping in mind that more instances means the cache is split more times. If you have a cache set to 500 but it has 10 instances, each cache will only have 50 cached.\\n\\nThe `table_definition_cache` and `table_open_cache` can be left as default as they are auto-sized MySQL 5.6 and above (ie: do not set them to any value).\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"percentUnit\",\"max\":1,\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"948ad10b-8b22-4d42-9e94-99ef09e12927\",\"layout\":{\"h\":7,\"i\":\"948ad10b-8b22-4d42-9e94-99ef09e12927\",\"isResizable\":true,\"w\":12,\"x\":12,\"y\":59},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"mysql_global_status_open_tables{asset_id=\\\"$asset_id\\\"}\",\"legend\":\"打开表数量(Open Tables)\"},{\"expr\":\"mysql_global_variables_table_open_cache{asset_id=\\\"$asset_id\\\"}\",\"legend\":\"打开缓存表数量(Table Open Cache)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"打开表数量(MySQL Open Tables)\",\"description\":\"**MySQL Open Tables**\\n\\nThe recommendation is to set the `table_open_cache_instances` to a loose correlation to virtual CPUs, keeping in mind that more instances means the cache is split more times. If you have a cache set to 500 but it has 10 instances, each cache will only have 50 cached.\\n\\nThe `table_definition_cache` and `table_open_cache` can be left as default as they are auto-sized MySQL 5.6 and above (ie: do not set them to any value).\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}}],\"var\":[{\"name\":\"prom\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"label\":\"\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${prom}\"},\"definition\":\"label_values(mysql_up, asset_id)\"}],\"version\":\"3.0.0\"}', 'MySQL');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (4, '{\"panels\":[{\"collapsed\":true,\"id\":\"2ecb82c6-4d1a-41b5-8cdc-0284db16bd54\",\"layout\":{\"h\":1,\"i\":\"2ecb82c6-4d1a-41b5-8cdc-0284db16bd54\",\"w\":24,\"x\":0,\"y\":0},\"name\":\"基本信息\",\"type\":\"row\"},{\"type\":\"stat\",\"id\":\"b5acc352-a2bd-4afc-b6cd-d6db0905f807\",\"layout\":{\"h\":3,\"i\":\"b5acc352-a2bd-4afc-b6cd-d6db0905f807\",\"w\":6,\"x\":0,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"min(redis_uptime_in_seconds{asset_id=\\\"$asset_id\\\"})/3600/24\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"启动时长(单位: 天)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{\"util\":\"none\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"8ccada5e-02f3-4efc-9b36-2a367612e4cb\",\"layout\":{\"h\":3,\"i\":\"8ccada5e-02f3-4efc-9b36-2a367612e4cb\",\"w\":6,\"x\":6,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(redis_connected_clients{asset_id=\\\"$asset_id\\\"})\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"客户端连接数(Connected Clients)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"716dc7e7-c9ec-4195-93f6-db1c572ae8b0\",\"layout\":{\"h\":3,\"i\":\"716dc7e7-c9ec-4195-93f6-db1c572ae8b0\",\"w\":6,\"x\":12,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"redis_used_memory{asset_id=\\\"$asset_id\\\"}\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"已使用内存(Memory Used)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"valueMappings\":[{\"match\":{\"from\":128000000},\"result\":{\"color\":\"#f10909\"},\"type\":\"range\"}],\"standardOptions\":{\"util\":\"bytesIEC\",\"decimals\":0},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"c6948161-db07-42df-beb1-765ee9c071a9\",\"layout\":{\"h\":3,\"i\":\"c6948161-db07-42df-beb1-765ee9c071a9\",\"w\":6,\"x\":18,\"y\":1},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"redis_maxmemory{asset_id=\\\"$asset_id\\\"}\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"最大内存(Max Memory Limit)\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{\"util\":\"bytesIEC\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"collapsed\":true,\"id\":\"bd54cf4f-1abb-4945-8aab-f89aec16daef\",\"layout\":{\"h\":1,\"i\":\"bd54cf4f-1abb-4945-8aab-f89aec16daef\",\"w\":24,\"x\":0,\"y\":4},\"name\":\"命令信息\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"3d5f8c4e-0ddf-4d68-9f6d-2cc57d864a8e\",\"layout\":{\"h\":7,\"i\":\"3d5f8c4e-0ddf-4d68-9f6d-2cc57d864a8e\",\"w\":8,\"x\":0,\"y\":5},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"rate(redis_total_commands_processed{asset_id=\\\"$asset_id\\\"}[5m])\",\"legend\":\"每秒执行\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"每秒命令执行(Commands Executed / sec)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"344a874d-c34d-4d2d-9bb4-46e0912cd9f5\",\"layout\":{\"h\":7,\"i\":\"344a874d-c34d-4d2d-9bb4-46e0912cd9f5\",\"w\":8,\"x\":8,\"y\":5},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"irate(redis_keyspace_hits{asset_id=\\\"$asset_id\\\"}[5m])\",\"legend\":\"命中(hits)\"},{\"expr\":\"irate(redis_keyspace_misses{asset_id=\\\"$asset_id\\\"}[5m])\",\"legend\":\"未命中(misses)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"命中率(Hits / Misses per Sec)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"3c83cd35-585c-4070-a210-1f17345f13f4\",\"layout\":{\"h\":7,\"i\":\"3c83cd35-585c-4070-a210-1f17345f13f4\",\"w\":8,\"x\":16,\"y\":5},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"topk(5, irate(redis_cmdstat_calls{asset_id=\\\"$asset_id\\\"} [1m]))\",\"legend\":\"{{command}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"命令分布(Top Commands)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"1ea61073-a46d-4d7c-b072-fcdcbc5ac084\",\"layout\":{\"h\":1,\"i\":\"1ea61073-a46d-4d7c-b072-fcdcbc5ac084\",\"w\":24,\"x\":0,\"y\":12},\"name\":\"数据监控(Keys)\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"b2b4451c-4f8a-438a-8c48-69c95c68361e\",\"layout\":{\"h\":7,\"i\":\"b2b4451c-4f8a-438a-8c48-69c95c68361e\",\"w\":8,\"x\":0,\"y\":13},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum (redis_keyspace_keys{asset_id=\\\"$asset_id\\\"}) by (db)\",\"legend\":\"{{db}}\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"数据量(Total Items per DB)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"894b9beb-e764-441c-ae04-13e5dbbb901d\",\"layout\":{\"h\":7,\"i\":\"894b9beb-e764-441c-ae04-13e5dbbb901d\",\"w\":8,\"x\":8,\"y\":13},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(redis_expired_keys{asset_id=\\\"$asset_id\\\"}[5m])) by (instance)\",\"legend\":\"过期(expired)\"},{\"expr\":\"sum(rate(redis_evicted_keys{asset_id=\\\"$asset_id\\\"}[5m])) by (instance)\",\"legend\":\"淘汰(evicted)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"过期/淘汰(Expired / Evicted)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"f721a641-28c7-4e82-a37c-ec17704a0c57\",\"layout\":{\"h\":7,\"i\":\"f721a641-28c7-4e82-a37c-ec17704a0c57\",\"w\":8,\"x\":16,\"y\":13},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(redis_keyspace_keys{asset_id=\\\"$asset_id\\\"}) - sum(redis_keyspace_expires{asset_id=\\\"$asset_id\\\"}) \",\"legend\":\"未过期(not expiring)\"},{\"expr\":\"sum(redis_keyspace_expires{asset_id=\\\"$asset_id\\\"}) \",\"legend\":\"过期(expiring)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"过期数据占比(Expiring vs Not-Expiring Keys)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"collapsed\":true,\"id\":\"60ff41ed-9d41-40ee-a13b-c968f3ca49d0\",\"layout\":{\"h\":1,\"i\":\"60ff41ed-9d41-40ee-a13b-c968f3ca49d0\",\"w\":24,\"x\":0,\"y\":20},\"name\":\"网络信息\",\"type\":\"row\"},{\"type\":\"timeseriesN\",\"id\":\"1841950c-e867-4a62-b846-78754dc0e34d\",\"layout\":{\"h\":7,\"i\":\"1841950c-e867-4a62-b846-78754dc0e34d\",\"w\":24,\"x\":0,\"y\":21},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"sum(rate(redis_total_net_input_bytes{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"入流量(input)\"},{\"expr\":\"sum(rate(redis_total_net_output_bytes{asset_id=\\\"$asset_id\\\"}[5m]))\",\"legend\":\"出流量(output)\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"网络流量(Network I/O)\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bytesIEC\",\"decimals\":1},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"row\",\"id\":\"c6ed879d-0f21-4270-8ac6-95e3ce3bfa1e\",\"name\":\"分组\",\"collapsed\":true,\"layout\":{\"x\":0,\"y\":28,\"w\":24,\"h\":1,\"i\":\"c6ed879d-0f21-4270-8ac6-95e3ce3bfa1e\"}},{\"type\":\"row\",\"id\":\"82eeec8b-512d-4250-9d09-a2948a719388\",\"name\":\"分组\",\"collapsed\":true,\"layout\":{\"x\":0,\"y\":29,\"w\":24,\"h\":1,\"i\":\"82eeec8b-512d-4250-9d09-a2948a719388\"}}],\"var\":[{\"name\":\"prom\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"label\":\"\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${prom}\"},\"definition\":\"label_values(redis_up, asset_id)\",\"multi\":false}],\"version\":\"3.0.0\"}', 'Redis');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (5, '{\"panels\":[{\"type\":\"table\",\"id\":\"3674dbfa-243a-49f6-baa5-b7f887c1afb0\",\"layout\":{\"h\":3,\"w\":24,\"x\":0,\"y\":0,\"i\":\"3674dbfa-243a-49f6-baa5-b7f887c1afb0\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${Datasource}\",\"targets\":[{\"expr\":\"max(http_response_result_code{asset_id=\\\"$asset_id\\\"}) by (target)\",\"instant\":true,\"legend\":\"状态\",\"refId\":\"A\"},{\"expr\":\"max(http_response_response_code{asset_id=\\\"$asset_id\\\"}) by (target)\",\"instant\":true,\"legend\":\"http状态码(status code)\",\"refId\":\"B\"},{\"expr\":\"max(http_response_response_time{asset_id=\\\"$asset_id\\\"}) by (target) *1000\",\"instant\":true,\"legend\":\"响应时间(latency)\",\"refId\":\"C\"},{\"expr\":\"max(http_response_cert_expire_timestamp{asset_id=\\\"$asset_id\\\"}) by (target) - time()\",\"instant\":true,\"legend\":\"证书有效期(cert expire)\",\"refId\":\"D\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"WEB服务监控(URL Details)\",\"maxPerRow\":4,\"custom\":{\"showHeader\":true,\"colorMode\":\"background\",\"calc\":\"lastNotNull\",\"displayMode\":\"labelValuesToRows\",\"sortColumn\":\"target\",\"sortOrder\":\"ascend\",\"aggrDimension\":\"target\"},\"options\":{\"valueMappings\":[],\"standardOptions\":{}},\"overrides\":[{\"matcher\":{\"value\":\"A\"},\"properties\":{\"standardOptions\":{},\"valueMappings\":[{\"match\":{\"special\":0},\"result\":{\"color\":\"#2c9d3d\",\"text\":\"在线\"},\"type\":\"special\"},{\"match\":{\"from\":1,\"special\":1},\"result\":{\"color\":\"#e90f0f\",\"text\":\"故障\"},\"type\":\"range\"}]}},{\"matcher\":{\"value\":\"D\"},\"properties\":{\"standardOptions\":{\"util\":\"humantimeSeconds\"},\"valueMappings\":[{\"match\":{\"to\":604800},\"result\":{\"color\":\"#f60c0c\"},\"type\":\"range\"},{\"match\":{\"to\":2592000},\"result\":{\"color\":\"#ffae39\"},\"type\":\"range\"},{\"type\":\"range\",\"result\":{\"color\":\"#2c9d3d\"},\"match\":{\"from\":2592000}}]},\"type\":\"special\"},{\"matcher\":{\"value\":\"B\"},\"properties\":{\"standardOptions\":{},\"valueMappings\":[{\"match\":{\"to\":399},\"result\":{\"color\":\"#2c9d3d\"},\"type\":\"range\"},{\"match\":{\"to\":499},\"result\":{\"color\":\"#ff656b\"},\"type\":\"range\"},{\"match\":{\"from\":500},\"result\":{\"color\":\"#f10808\"},\"type\":\"range\"}]},\"type\":\"special\"},{\"matcher\":{\"value\":\"C\"},\"properties\":{\"standardOptions\":{\"util\":\"milliseconds\"},\"valueMappings\":[{\"match\":{\"to\":400},\"result\":{\"color\":\"#2c9d3d\"},\"type\":\"range\"},{\"match\":{\"from\":400},\"result\":{\"color\":\"#ff656b\"},\"type\":\"range\"},{\"match\":{\"from\":2000},\"result\":{\"color\":\"#f11313\"},\"type\":\"range\"}]},\"type\":\"special\"}]},{\"type\":\"timeseriesN\",\"id\":\"b996687f-a54d-46e1-8f07-c554c4e2bf49\",\"layout\":{\"h\":5,\"w\":24,\"x\":0,\"y\":3,\"i\":\"b996687f-a54d-46e1-8f07-c554c4e2bf49\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${Datasource}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"http_response_response_time{asset_id=\\\"$asset_id\\\"} * 1000\",\"legend\":\"响应时间\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"响应时间\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"milliseconds\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}}],\"var\":[{\"name\":\"Datasource\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"label\":\"\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${Datasource}\"},\"definition\":\"label_values(http_response_result_code, asset_id)\"}],\"version\":\"3.0.0\"}', 'HTTP服务');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (6, '{\"panels\":[{\"type\":\"table\",\"id\":\"73c6eaf9-1685-4a7a-bf53-3d52afa1792e\",\"layout\":{\"h\":3,\"w\":24,\"x\":0,\"y\":0,\"i\":\"73c6eaf9-1685-4a7a-bf53-3d52afa1792e\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"max(net_response_result_code{asset_id=\\\"${asset_id}\\\"}) by (target)\",\"legend\":\"状态\",\"refId\":\"A\"},{\"expr\":\"max(net_response_response_time{asset_id=\\\"${asset_id}\\\"}) by (target) * 1000\",\"legend\":\"响应时长(ms)\",\"refId\":\"C\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{\"indexByName\":{\"target\":0}}}],\"name\":\"端口检测\",\"maxPerRow\":4,\"custom\":{\"showHeader\":true,\"colorMode\":\"background\",\"calc\":\"lastNotNull\",\"displayMode\":\"labelValuesToRows\",\"aggrDimension\":\"target\"},\"options\":{\"valueMappings\":[],\"standardOptions\":{}},\"overrides\":[{\"matcher\":{\"value\":\"A\"},\"properties\":{\"standardOptions\":{},\"valueMappings\":[{\"match\":{\"special\":0},\"result\":{\"color\":\"#2c9d3d\",\"text\":\"在线\"},\"type\":\"special\"},{\"match\":{\"from\":1,\"special\":1},\"result\":{\"color\":\"#e90f0f\",\"text\":\"故障\"},\"type\":\"range\"}]}},{\"type\":\"special\",\"matcher\":{\"value\":\"C\"},\"properties\":{\"valueMappings\":[{\"type\":\"range\",\"result\":{\"color\":\"#f10c0c\"},\"match\":{\"from\":1}},{\"type\":\"range\",\"result\":{\"color\":\"#2c9d3d\"},\"match\":{\"to\":1}}],\"standardOptions\":{\"util\":\"milliseconds\",\"decimals\":3}}}]},{\"type\":\"timeseriesN\",\"id\":\"1dcb028f-492a-4b7a-9c65-c8847f059d15\",\"layout\":{\"h\":5,\"w\":24,\"x\":0,\"y\":3,\"i\":\"1dcb028f-492a-4b7a-9c65-c8847f059d15\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"net_response_response_time{asset_id=\\\"${asset_id}\\\"} * 1000\",\"legend\":\"响应时长\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"响应时长\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"milliseconds\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}}],\"var\":[{\"name\":\"prom\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${prom}\"},\"definition\":\"label_values(net_response_result_code,asset_id)\"}],\"version\":\"3.0.0\"}', '网络服务');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (7, '{\"panels\":[{\"type\":\"table\",\"id\":\"1677138f-0f33-485c-8ee1-2db24cabbf54\",\"layout\":{\"h\":3,\"w\":24,\"x\":0,\"y\":0,\"i\":\"1677138f-0f33-485c-8ee1-2db24cabbf54\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"max(ping_result_code) by (target)\",\"legend\":\"状态\",\"refId\":\"A\"},{\"expr\":\"max(ping_percent_packet_loss) by (target)\",\"legend\":\"丢包率(%)\",\"refId\":\"B\"},{\"expr\":\"max(ping_maximum_response_ms) by (target) \",\"legend\":\"响应时间(ms)\",\"refId\":\"C\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"Ping\",\"maxPerRow\":4,\"custom\":{\"showHeader\":true,\"colorMode\":\"background\",\"calc\":\"lastNotNull\",\"displayMode\":\"labelValuesToRows\",\"aggrDimension\":\"target\"},\"options\":{\"valueMappings\":[],\"standardOptions\":{}},\"overrides\":[{\"matcher\":{\"value\":\"A\"},\"properties\":{\"standardOptions\":{},\"valueMappings\":[{\"match\":{\"special\":0},\"result\":{\"color\":\"#2c9d3d\",\"text\":\"正常\"},\"type\":\"special\"},{\"match\":{\"from\":1,\"special\":1},\"result\":{\"color\":\"#e90f0f\",\"text\":\"故障\"},\"type\":\"range\"}]}},{\"type\":\"special\",\"matcher\":{\"value\":\"B\"},\"properties\":{\"valueMappings\":[{\"type\":\"range\",\"result\":{\"color\":\"#f30a0a\"},\"match\":{\"from\":1}},{\"type\":\"special\",\"result\":{\"color\":\"#2c9d3d\"},\"match\":{\"special\":0}}],\"standardOptions\":{\"util\":\"percent\",\"decimals\":1}}},{\"type\":\"special\",\"matcher\":{\"value\":\"C\"},\"properties\":{\"valueMappings\":[{\"type\":\"range\",\"result\":{\"color\":\"#2c9d3d\"},\"match\":{\"from\":null,\"to\":100}},{\"type\":\"range\",\"result\":{\"color\":\"#ff8286\"},\"match\":{\"to\":300}},{\"type\":\"range\",\"result\":{\"color\":\"#f00808\"},\"match\":{\"to\":null,\"from\":1000}}],\"standardOptions\":{\"util\":\"milliseconds\"}}}]},{\"type\":\"timeseriesN\",\"id\":\"fdbbe917-0a8d-4190-9005-c052fa24a10e\",\"layout\":{\"h\":5,\"w\":12,\"x\":0,\"y\":3,\"i\":\"fdbbe917-0a8d-4190-9005-c052fa24a10e\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"ping_average_response_ms\",\"legend\":\"{{target}}-响应时间\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"响应时间\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"milliseconds\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"35d05420-55c2-42ac-84b0-da0235852f79\",\"layout\":{\"h\":5,\"w\":12,\"x\":12,\"y\":3,\"i\":\"9a2d68c5-8891-40b2-bc52-9930c8873974\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"ping_percent_packet_loss\",\"legend\":\"{{target}}-丢包率\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"丢包率\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"percent\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}}],\"var\":[{\"name\":\"prom\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${prom}\"},\"definition\":\"label_values(ping_result_code, asset_id)\"}],\"version\":\"3.0.0\"}', '网络端点');
INSERT INTO `board_payload`(`id`, `payload`, `asset_type`) VALUES (8, '{\"links\":[],\"panels\":[{\"type\":\"stat\",\"id\":\"d5e905cf-da22-48be-9fca-1f92695ca730\",\"layout\":{\"h\":4,\"w\":6,\"x\":0,\"y\":0,\"i\":\"d5e905cf-da22-48be-9fca-1f92695ca730\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"switch_legacy_uptime{asset_id=\\\"$asset_id\\\"}/100/3600/24\",\"instant\":true,\"legend\":\"\",\"refId\":\"A\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"已运行时间(单位: 天)\",\"links\":[],\"description\":\"系统启动时间\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"graphMode\":\"none\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{\"util\":\"none\",\"decimals\":0},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"gaugeN\",\"id\":\"d755c99d-a323-41e6-8117-6bc006bef8b7\",\"layout\":{\"h\":4,\"w\":6,\"x\":6,\"y\":0,\"i\":\"bd2cd5b0-50ac-42d7-b29d-ea89ceb015a7\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"switch_legacy_cpu_util{asset_id=\\\"$asset_id\\\"}\",\"legend\":\"\",\"refId\":\"A\",\"instant\":false}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"CPU 使用率 %\",\"links\":[],\"description\":\"\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{\"util\":\"percent\"},\"thresholds\":{\"steps\":[{\"color\":\"#3FC453\",\"value\":null,\"type\":\"base\"},{\"color\":\"#FF9919\",\"value\":60},{\"color\":\"#FF656B\",\"value\":80}]}}},{\"type\":\"gaugeN\",\"id\":\"c3991b49-1ad8-4f63-87b8-d41bbf729833\",\"layout\":{\"h\":4,\"w\":6,\"x\":12,\"y\":0,\"i\":\"109aad94-79bd-4aec-b8ac-db73cb6601a8\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"switch_legacy_mem_util{asset_id=\\\"$asset_id\\\"}\",\"legend\":\"mem_usage\",\"refId\":\"A\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"内存使用率 %\",\"links\":[],\"description\":\"内存使用率 \",\"maxPerRow\":4,\"custom\":{\"textMode\":\"value\",\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{\"util\":\"percent\"},\"thresholds\":{\"steps\":[{\"color\":\"#3FC453\",\"value\":null,\"type\":\"base\"},{\"color\":\"#FF9919\",\"value\":60},{\"color\":\"#FF656B\",\"value\":80}]}}},{\"type\":\"gaugeN\",\"id\":\"616de58a-70a7-4c0b-b0f2-5151b9f0e9c5\",\"layout\":{\"h\":4,\"w\":6,\"x\":18,\"y\":0,\"i\":\"616de58a-70a7-4c0b-b0f2-5151b9f0e9c5\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"topk(1, switch_legacy_if_in_speed_percent{asset_id=\\\"$asset_id\\\"} and switch_legacy_if_out_speed_percent{asset_id=\\\"$asset_id\\\"})\",\"instant\":true,\"legend\":\"{{ifname}}\",\"refId\":\"A\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"最大端口带宽利用率\",\"links\":[],\"description\":\"交换机有多个端口, 此处显示带宽使用最大的端口的利用率\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"valueAndName\",\"calc\":\"lastNotNull\"},\"options\":{\"standardOptions\":{\"util\":\"percent\"},\"thresholds\":{\"steps\":[{\"color\":\"#3FC453\",\"value\":null,\"type\":\"base\"},{\"color\":\"#FF9919\",\"value\":60},{\"color\":\"#FF656B\",\"value\":80}]}}},{\"type\":\"hexbin\",\"id\":\"d909ec3d-3b1f-4156-9b43-f1cdf1e549e1\",\"layout\":{\"h\":4,\"w\":12,\"x\":0,\"y\":4,\"i\":\"d909ec3d-3b1f-4156-9b43-f1cdf1e549e1\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"switch_legacy_if_oper_status{asset_id=\\\"${asset_id}\\\"}\",\"legend\":\"{{ifname}}\",\"instant\":true}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"端口状态\",\"description\":\"绿色为Up状态, 红色为Down状态\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"name\",\"calc\":\"lastNotNull\",\"colorRange\":[\"#83c898\",\"#c2c2c2\",\"#fc653f\"],\"reverseColorOrder\":false,\"colorDomainAuto\":true},\"options\":{\"standardOptions\":{}}},{\"type\":\"stat\",\"id\":\"5de5d2d2-e565-4dd1-bc29-0fefade35c62\",\"layout\":{\"h\":4,\"w\":6,\"x\":12,\"y\":4,\"i\":\"5de5d2d2-e565-4dd1-bc29-0fefade35c62\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"sum(switch_legacy_if_in_discards{asset_id=\\\"${asset_id}\\\"} and switch_legacy_if_out_discards{asset_id=\\\"${asset_id}\\\"})\",\"legend\":\"\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"丢包数量\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"valueAndName\",\"graphMode\":\"area\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#ff9919\",\"value\":1,\"type\":\"\"},{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"stat\",\"id\":\"d454150e-e2f2-4e2b-bf65-62ecf03fcdf1\",\"layout\":{\"h\":4,\"w\":6,\"x\":18,\"y\":4,\"i\":\"abdac050-5363-4baf-b4a2-58697cf63743\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"refId\":\"A\",\"expr\":\"sum(switch_legacy_if_in_errors{asset_id=\\\"${asset_id}\\\"} and switch_legacy_if_out_errors{asset_id=\\\"${asset_id}\\\"})\",\"legend\":\"\"}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"错包数量\",\"maxPerRow\":4,\"custom\":{\"textMode\":\"valueAndName\",\"graphMode\":\"area\",\"colorMode\":\"value\",\"calc\":\"lastNotNull\",\"valueField\":\"Value\",\"colSpan\":1,\"textSize\":{}},\"options\":{\"standardOptions\":{},\"thresholds\":{\"steps\":[{\"color\":\"#ff9919\",\"value\":1,\"type\":\"\"},{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}}},{\"type\":\"timeseriesN\",\"id\":\"26ae7fc1-230e-451e-9415-ea93ae8b2abb\",\"layout\":{\"h\":4,\"w\":12,\"x\":0,\"y\":8,\"i\":\"26ae7fc1-230e-451e-9415-ea93ae8b2abb\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"switch_legacy_if_in{asset_id=\\\"$asset_id\\\"} > 0\",\"instant\":false,\"legend\":\"{{ifname}}\",\"refId\":\"A\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"端口入流量\",\"links\":[],\"description\":\"为0的不显示\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bitsIEC\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}},{\"type\":\"timeseriesN\",\"id\":\"15cd3b8d-dcb9-4e87-833e-003c93fc12bf\",\"layout\":{\"h\":4,\"w\":12,\"x\":12,\"y\":8,\"i\":\"ae973664-ec5d-45fd-ac8c-a644b843f55f\",\"isResizable\":true},\"version\":\"3.0.0\",\"datasourceCate\":\"prometheus\",\"datasourceValue\":\"${prom}\",\"targets\":[{\"expr\":\"switch_legacy_if_out{asset_id=\\\"$asset_id\\\"} > 0\",\"instant\":false,\"legend\":\"{{ifname}}\",\"refId\":\"A\",\"step\":null}],\"transformations\":[{\"id\":\"organize\",\"options\":{}}],\"name\":\"端口出流量\",\"links\":[],\"description\":\"为0的不显示\",\"maxPerRow\":4,\"options\":{\"tooltip\":{\"mode\":\"all\",\"sort\":\"none\"},\"legend\":{\"displayMode\":\"hidden\"},\"standardOptions\":{\"util\":\"bitsIEC\"},\"thresholds\":{\"steps\":[{\"color\":\"#6395fa\",\"value\":null,\"type\":\"base\"}]}},\"custom\":{\"lineInterpolation\":\"smooth\",\"fillOpacity\":0.5,\"stack\":\"off\"}}],\"var\":[{\"name\":\"prom\",\"label\":\"数据源\",\"type\":\"datasource\",\"hide\":true,\"definition\":\"prometheus\"},{\"name\":\"asset_id\",\"label\":\"\",\"type\":\"query\",\"hide\":true,\"datasource\":{\"cate\":\"prometheus\",\"value\":\"${prom}\"},\"definition\":\"label_values(switch_legacy_uptime, asset_id)\"}],\"version\":\"3.0.0\"}', '网络设备');


-- deprecated
CREATE TABLE `dashboard` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default 0 comment 'busi group id',
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'split by space',
    `configs` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci comment 'dashboard variables',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`group_id`, `name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- deprecated
-- auto create the first subclass 'Default chart group' of dashboard
CREATE TABLE `chart_group` (
    `id` bigint unsigned not null auto_increment,
    `dashboard_id` bigint unsigned not null,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `weight` int not null default 0,
    PRIMARY KEY (`id`),
    KEY (`dashboard_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- deprecated
CREATE TABLE `chart` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint unsigned not null comment 'chart group id',
    `configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `weight` int not null default 0,
    PRIMARY KEY (`id`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `chart_share` (
    `id` bigint unsigned not null auto_increment,
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_id` bigint unsigned not null default 0,
    `configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    primary key (`id`),
    key (`create_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `alert_rule` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default 0 comment 'busi group id',
    `asset_id` bigint DEFAULT NULL COMMENT '资产id',
    `asset_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产名称',
    `asset_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产IP',
    `cate` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'datasource ids',
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `note` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `prod` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `algorithm` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `algo_params` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `delay` int not null default 0,
    `severity` tinyint(1) not null comment '1:Emergency 2:Warning 3:Notice',
    `disabled` tinyint(1) not null comment '0:enabled 1:disabled',
    `prom_for_duration` int not null comment 'prometheus for, unit:s',
    `rule_config_cn` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'rule_config_cn',
    `rule_config` text  CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'rule_config',
    `rule_config_fe` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci comment 'rule_config',
    `prom_ql` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'promql',
    `prom_eval_interval` int not null comment 'evaluate interval',
    `enable_stime` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '00:00',
    `enable_etime` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '23:59',
    `enable_days_of_week` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: 0 1 2 3 4 5 6',
    `enable_in_bg` tinyint(1) not null default 0 comment '1: only this bg 0: global',
    `notify_recovered` tinyint(1) not null comment 'whether notify when recovery',
    `notify_channels` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: sms voice email dingtalk wecom',
    `notify_groups` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: 233 43',
    `notify_repeat_step` int not null default 0 comment 'unit: min',
    `notify_max_number` int not null default 0 comment '',
    `recover_duration` int not null default 0 comment 'unit: s',
    `callbacks` varchar(1024) not null default '' comment 'split by space: http://a.com/api/x http://a.com/api/y',
    `runbook_url` varchar(255),
    `append_tags` varchar(255) not null default '' comment 'split by space: service=n9e mod=api',
    `annotations` text not null comment 'annotations',
    `extra_config` text not null comment 'extra_config',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `deleted_at` datetime(0) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY (`group_id`),
    KEY (`update_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `alert_mute` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default 0 comment 'busi group id',
    `prod` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `note` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `cate` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'datasource ids',
    `tags` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'json,map,tagkey->regexp|value',
    `cause` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `btime` bigint not null default 0 comment 'begin time',
    `etime` bigint not null default 0 comment 'end time',
    `disabled` tinyint(1) not null default 0 comment '0:enabled 1:disabled',
    `mute_time_type` tinyint(1) not null default 0,
    `periodic_mutes` varchar(4096) not null default '',
    `severities` varchar(32) not null default '',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`create_at`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `alert_subscribe` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `disabled` tinyint(1) not null default 0 comment '0:enabled 1:disabled',
    `group_id` bigint not null default 0 comment 'busi group id',
    `prod` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `cate` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'datasource ids',
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `rule_id` bigint not null default 0,
    `severities` varchar(32) not null default '',
    `tags` varchar(4096) not null default '' comment 'json,map,tagkey->regexp|value',
    `redefine_severity` tinyint(1) default 0 comment 'is redefine severity?',
    `new_severity` tinyint(1) not null comment '0:Emergency 1:Warning 2:Notice',
    `redefine_channels` tinyint(1) default 0 comment 'is redefine channels?',
    `new_channels` varchar(255) not null default '' comment 'split by space: sms voice email dingtalk wecom',
    `user_group_ids` varchar(250) not null comment 'split by space 1 34 5, notify cc to user_group_ids',
    `webhooks` text not null,
    `extra_config` text not null comment 'extra_config',
    `redefine_webhooks` tinyint(1) default 0,
    `for_duration` bigint not null default 0,
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`update_at`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `target` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default 0 comment 'busi group id',
    `ident` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'target id',
    `current_version` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'append to alert event as field',
    `tags` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'append to series data as tags, split by space, append external space at suffix',
    `update_at` bigint not null default 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`ident`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;


-- case1: target_idents; case2: target_tags
-- CREATE TABLE `collect_rule` (
--     `id` bigint unsigned not null auto_increment,
--     `group_id` bigint not null default 0 comment 'busi group id',
--     `cluster` varchar(128) not null,
--     `target_idents` varchar(512) not null default '' comment 'ident list, split by space',
--     `target_tags` varchar(512) not null default '' comment 'filter targets by tags, split by space',
--     `name` varchar(191) not null default '',
--     `note` varchar(255) not null default '',
--     `step` int not null,
--     `type` varchar(64) not null comment 'e.g. port proc log plugin',
--     `data` text not null,
--     `append_tags` varchar(255) not null default '' comment 'split by space: e.g. mod=n9e dept=cloud',
--     `create_at` bigint not null default 0,
--     `create_by` varchar(64) not null default '',
--     `update_at` bigint not null default 0,
--     `update_by` varchar(64) not null default '',
--     PRIMARY KEY (`id`),
--     KEY (`group_id`, `type`, `name`)
-- ) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `metric_view` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `cate` tinyint(1) not null comment '0: preset 1: custom',
    `configs` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `create_at` bigint not null default 0,
    `create_by` bigint not null default 0 comment 'user id',
    `update_at` bigint not null default 0,
    PRIMARY KEY (`id`),
    KEY (`create_by`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into metric_view(name, cate, configs) values('Host View', 0, '{"filters":[{"oper":"=","label":"__name__","value":"cpu_usage_idle"}],"dynamicLabels":[],"dimensionLabels":[{"label":"ident","value":""}]}');

CREATE TABLE `recording_rule` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default '0' comment 'group_id',
    `datasource_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'datasource ids',
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'new metric name',
    `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'rule note',
    `disabled` tinyint(1) not null default 0 comment '0:enabled 1:disabled',
    `prom_ql` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'promql',
    `prom_eval_interval` int not null comment 'evaluate interval',
    `append_tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '' comment 'split by space: service=n9e mod=api',
    `query_configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'query configs',
    `create_at` bigint default '0',
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '',
    `update_at` bigint default '0',
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '',
    PRIMARY KEY (`id`),
    KEY `group_id` (`group_id`),
    KEY `update_at` (`update_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `alert_aggr_view` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `rule` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `cate` tinyint(1) not null comment '0: preset 1: custom',
    `create_at` bigint not null default 0,
    `create_by` bigint not null default 0 comment 'user id',
    `update_at` bigint not null default 0,
    PRIMARY KEY (`id`),
    KEY (`create_by`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

insert into alert_aggr_view(name, rule, cate) values('By BusiGroup, Severity', 'field:group_name::field:severity', 0);
insert into alert_aggr_view(name, rule, cate) values('By RuleName', 'field:rule_name', 0);

CREATE TABLE `alert_cur_event` (
    `id` bigint unsigned not null comment 'use alert_his_event.id',
    `asset_id` bigint DEFAULT NULL COMMENT '资产id',
    `asset_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产名称',
    `asset_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产IP',
    `cate` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_id` bigint not null default 0 comment 'datasource id',
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `group_id` bigint unsigned not null comment 'busi group id of rule',
    `group_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'busi group name',
    `hash` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'rule_id + vector_pk',
    `rule_id` bigint unsigned not null,
    `rule_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `rule_note` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default 'alert rule note',
    `rule_prod` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `rule_algo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `severity` tinyint(1) not null comment '0:Emergency 1:Warning 2:Notice',
    `prom_for_duration` int not null comment 'prometheus for, unit:s',
    `prom_ql` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'promql',
    `prom_eval_interval` int not null comment 'evaluate interval',
    `callbacks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: http://a.com/api/x http://a.com/api/y',
    `runbook_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `notify_recovered` tinyint(1) not null comment 'whether notify when recovery',
    `notify_channels` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: sms voice email dingtalk wecom',
    `notify_groups` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: 233 43',
    `notify_repeat_next` bigint not null default 0 comment 'next timestamp to notify, get repeat settings from rule',
    `notify_cur_number` int not null default 0 comment '',
    `target_ident` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'target ident, also in tags',
    `target_note` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'target note',
    `first_trigger_time` bigint,
    `trigger_time` bigint not null,
    `trigger_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `annotations` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'annotations',
    `rule_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'annotations',
    `tags` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'merge data_tags rule_tags, split by ,,',
    `status` tinyint(1) not null default 0 comment '状态',
    `deleted_at` datetime(0) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY (`hash`),
    KEY (`rule_id`),
    KEY (`trigger_time`, `group_id`),
    KEY (`notify_repeat_next`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `alert_his_event` (
    `id` bigint unsigned not null AUTO_INCREMENT,
    `is_recovered` tinyint(1) not null,
    `asset_id` bigint DEFAULT NULL COMMENT '资产id',
    `asset_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产名称',
    `asset_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产IP',
    `cate` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_id` bigint not null default 0 comment 'datasource id',
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `group_id` bigint unsigned not null comment 'busi group id of rule',
    `group_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'busi group name',
    `hash` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'rule_id + vector_pk',
    `rule_id` bigint unsigned not null,
    `rule_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `rule_note` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default 'alert rule note',
    `rule_prod` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `rule_algo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `severity` tinyint(1) not null comment '0:Emergency 1:Warning 2:Notice',
    `prom_for_duration` int not null comment 'prometheus for, unit:s',
    `prom_ql` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'promql',
    `prom_eval_interval` int not null comment 'evaluate interval',
    `callbacks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: http://a.com/api/x http://a.com/api/y',
    `runbook_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `notify_recovered` tinyint(1) not null comment 'whether notify when recovery',
    `notify_channels` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: sms voice email dingtalk wecom',
    `notify_groups` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: 233 43',
    `notify_cur_number` int not null default 0 comment '',
    `target_ident` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'target ident, also in tags',
    `target_note` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'target note',
    `first_trigger_time` bigint,
    `trigger_time` bigint not null,
    `trigger_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `recover_time` bigint not null default 0,
    `last_eval_time` bigint not null default 0 comment 'for time filter',
    `tags` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'merge data_tags rule_tags, split by ,,',
    `annotations` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'annotations',
    `rule_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'annotations',
    `status` tinyint(1) not null default 0 comment '状态',
    `handle_by` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `handle_at` bigint not null default 0,
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `deleted_at` datetime(0) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY (`hash`),
    KEY (`rule_id`),
    KEY (`trigger_time`, `group_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `task_tpl`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `group_id`  int unsigned not null comment 'busi group id',
    `title`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `account`   varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  not null,
    `batch`     int unsigned not null default 0,
    `tolerance` int unsigned not null default 0,
    `timeout`   int unsigned not null default 0,
    `pause`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `script`    text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci         not null,
    `args`      varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `tags`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `task_tpl_host`
(
    `ii`   int unsigned NOT NULL AUTO_INCREMENT,
    `id`   int unsigned not null comment 'task tpl id',
    `host` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  not null comment 'ip or hostname',
    PRIMARY KEY (`ii`),
    KEY (`id`, `host`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `task_record`
(
    `id` bigint unsigned not null comment 'ibex task id',
    `event_id` bigint not null comment 'event id' default 0,
    `group_id` bigint not null comment 'busi group id',
    `ibex_address`   varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `ibex_auth_user` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `ibex_auth_pass` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `title`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci    not null default '',
    `account`   varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci     not null,
    `batch`     int unsigned    not null default 0,
    `tolerance` int unsigned    not null default 0,
    `timeout`   int unsigned    not null default 0,
    `pause`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci    not null default '',
    `script`    text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci            not null,
    `args`      varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci    not null default '',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`create_at`, `group_id`),
    KEY (`create_by`),
    KEY (`event_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `alerting_engines`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `instance` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'instance identification, e.g. 10.9.0.9:9090',
    `datasource_id` bigint not null default 0 comment 'datasource id',
    `engine_cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'n9e-alert cluster',
    `clock` bigint not null,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `datasource`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `category` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `plugin_id` int unsigned not null default 0,
    `plugin_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `plugin_type_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `cluster_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `settings` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `http` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `auth` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `created_at` bigint not null default 0,
    `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `updated_at` bigint not null default 0,
    `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    UNIQUE KEY (`name`),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `builtin_cate` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `user_id` bigint not null default 0,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `notify_tpl` (
    `id` bigint unsigned not null auto_increment,
    `channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`channel`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

INSERT INTO `notify_tpl` VALUES (1, 'Email', '邮件', '<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n    <title>一体化综合运维管理平台-告警通知</title>\n    <style type=\"text/css\">\n        .wrapper {\n            background-color: #f8f8f8;\n            padding: 15px;\n            height: 100%;\n        }\n        .main {\n            width: 600px;\n            padding: 30px;\n            margin: 0 auto;\n            background-color: #fff;\n            font-size: 12px;\n            font-family: verdana,\'Microsoft YaHei\',Consolas,\'Deja Vu Sans Mono\',\'Bitstream Vera Sans Mono\';\n        }\n        header {\n            border-radius: 2px 2px 0 0;\n        }\n        header .title {\n            font-size: 14px;\n            color: #333333;\n            margin: 0;\n        }\n        header .sub-desc {\n            color: #333;\n            font-size: 14px;\n            margin-top: 6px;\n            margin-bottom: 0;\n        }\n        hr {\n            margin: 20px 0;\n            height: 0;\n            border: none;\n            border-top: 1px solid #e5e5e5;\n        }\n        em {\n            font-weight: 600;\n        }\n        table {\n            margin: 20px 0;\n            width: 100%;\n        }\n\n        table tbody tr{\n            font-weight: 200;\n            font-size: 12px;\n            color: #666;\n            height: 32px;\n        }\n\n        .succ {\n            background-color: green;\n            color: #fff;\n        }\n\n        .fail {\n            background-color: red;\n            color: #fff;\n        }\n\n        .succ th, .succ td, .fail th, .fail td {\n            color: #fff;\n        }\n\n        table tbody tr th {\n            width: 80px;\n            text-align: right;\n        }\n        .text-right {\n            text-align: right;\n        }\n        .body {\n            margin-top: 24px;\n        }\n        .body-text {\n            color: #666666;\n            -webkit-font-smoothing: antialiased;\n        }\n        .body-extra {\n            -webkit-font-smoothing: antialiased;\n        }\n        .body-extra.text-right a {\n            text-decoration: none;\n            color: #333;\n        }\n        .body-extra.text-right a:hover {\n            color: #666;\n        }\n        .button {\n            width: 200px;\n            height: 50px;\n            margin-top: 20px;\n            text-align: center;\n            border-radius: 2px;\n            background: #2D77EE;\n            line-height: 50px;\n            font-size: 20px;\n            color: #FFFFFF;\n            cursor: pointer;\n        }\n        .button:hover {\n            background: rgb(25, 115, 255);\n            border-color: rgb(25, 115, 255);\n            color: #fff;\n        }\n        footer {\n            margin-top: 10px;\n            text-align: right;\n        }\n        .footer-logo {\n            text-align: right;\n        }\n        .footer-logo-image {\n            width: 108px;\n            height: 27px;\n            margin-right: 10px;\n        }\n        .copyright {\n            margin-top: 10px;\n            font-size: 12px;\n            text-align: right;\n            color: #999;\n            -webkit-font-smoothing: antialiased;\n        }\n    </style>\n</head>\n<body>\n<div class=\"wrapper\">\n    <div class=\"main\">\n        <header>\n            <h3 class=\"title\">{{.RuleName}}</h3>\n            <p class=\"sub-desc\"></p>\n        </header>\n\n        <hr>\n\n        <div class=\"body\">\n            <table cellspacing=\"0\" cellpadding=\"0\" border=\"0\">\n                <tbody>\n                {{if .IsRecovered}}\n                <tr class=\"succ\">\n                    <th>级别状态：</th>\n                    <td>S{{.Severity}} Recovered</td>\n                </tr>\n                {{else}}\n                <tr class=\"fail\">\n                    <th>级别状态：</th>\n                    <td>S{{.Severity}} Triggered</td>\n                </tr>\n                {{end}}\n\n                <tr>\n                    <th>策略备注：</th>\n                    <td>{{.RuleNote}}</td>\n                </tr>\n                <tr>\n                    <th>设备备注：</th>\n                    <td>{{.TargetNote}}</td>\n                </tr>\n                {{if not .IsRecovered}}\n                <tr>\n                    <th>触发时值：</th>\n                    <td>{{.TriggerValue}}</td>\n                </tr>\n                {{end}}\n\n                {{if .TargetIdent}}\n                <tr>\n                    <th>监控对象：</th>\n                    <td>{{.TargetIdent}}</td>\n                </tr>\n                {{end}}\n                <tr>\n                    <th>监控指标：</th>\n                    <td>{{.TagsJSON}}</td>\n                </tr>\n\n                {{if .IsRecovered}}\n                <tr>\n                    <th>恢复时间：</th>\n                    <td>{{timeformat .LastEvalTime}}</td>\n                </tr>\n                {{else}}\n                <tr>\n                    <th>触发时间：</th>\n                    <td>\n                        {{timeformat .TriggerTime}}\n                    </td>\n                </tr>\n                {{end}}\n\n                <tr>\n                    <th>发送时间：</th>\n                    <td>\n                        {{timestamp}}\n                    </td>\n                </tr>\n                </tbody>\n            </table>\n\n            <hr>\n\n            <footer>\n                <div class=\"copyright\" style=\"font-style: italic\">\n                    我们希望与您一起，将监控这个事情，做到极致！\n                </div>\n            </footer>\n        </div>\n    </div>\n</div>\n</body>\n</html>');


CREATE TABLE `sso_config` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `assets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `ident` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `group_id` bigint(20) NOT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `ip` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'IP地址',
  `manufacturers` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '厂商',
  `position` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产位置',
  `memo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `tags` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `plugin` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `label` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `params` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` int(0) NOT NULL DEFAULT 0,
  `status_at` bigint(0) NOT NULL DEFAULT 0,
  `directory_id` bigint(0) DEFAULT NULL COMMENT '所在分组',
  `payload` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
  `create_at` bigint(0) NOT NULL DEFAULT 0,
  `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `update_at` bigint(0) NOT NULL DEFAULT 0,
  `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `organization_id` bigint(0) DEFAULT NULL,
  `optional_metrics` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `deleted_at` datetime(0) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `group_id` (`group_id`),
  KEY `ident` (`ident`),
  KEY `organization_id` (`organization_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `organization`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `parent_id` int(0) DEFAULT NULL,
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `son` int(0) DEFAULT NULL,
  `city` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `manger` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `create_at` bigint(0) NOT NULL DEFAULT 0,
  `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `update_at` bigint(0) NOT NULL DEFAULT 0,
  `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `deleted_at` datetime(0) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE `es_index_pattern` (
    `id` bigint unsigned not null auto_increment,
    `datasource_id` bigint not null default 0 comment 'datasource id',
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `time_field` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '@timestamp',
    `allow_hide_system_indices` tinyint(1) not null default 0,
    `fields_format` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `create_at` bigint default '0',
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '',
    `update_at` bigint default '0',
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`datasource_id`, `name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

INSERT INTO `datasource`(`id`, `name`, `description`, `category`, `plugin_id`, `plugin_type`, `plugin_type_name`, `cluster_name`, `settings`, `status`, `http`, `auth`) VALUES (1, 'default', '', '', 0, 'prometheus', '', 'default', '{\"write_addr\":\"http://127.0.0.1:8428/api/v1/write\"}', 'enabled', '{\"timeout\":10000,\"dial_timeout\":0,\"tls\":{\"skip_tls_verify\":false},\"max_idle_conns_per_host\":0,\"url\":\"http://127.0.0.1:8428\",\"headers\":{}}', '{\"basic_auth\":false,\"basic_auth_user\":\"vm\",\"basic_auth_password\":\"vmdctbcab\"}');


-- ----------------------------
-- Table structure for alert_inspection_schedule
-- ----------------------------
DROP TABLE IF EXISTS `alert_inspection_schedule`;
CREATE TABLE `alert_inspection_schedule`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `paln_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `supervisor` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `description` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `area` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `report` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `time` bigint(0) DEFAULT NULL,
  `receiver` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `state` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `handle_by` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `handle_at` bigint(0) DEFAULT NULL,
  `update_at` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `update_by` bigint(0) DEFAULT NULL,
  `reset` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for asset_alter
-- ----------------------------
DROP TABLE IF EXISTS `asset_alter`;
CREATE TABLE `asset_alter`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSET_ID` int(0) NOT NULL COMMENT '资产ID',
  `ALTER_AT` int(0) DEFAULT NULL COMMENT '变更日期',
  `ALTER_EVENT_CODE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '变更事项编码',
  `ALTER_EVENT_KEY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '变更事项标签',
  `BEFORE_ALTER` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '变更前',
  `AFTER_ALTER` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '变更后',
  `ALTER_SPONSOR` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '变更发起人',
  `ALTER_STATUS` int(0) DEFAULT NULL COMMENT '确认状态(0:未确认;1:确认)',
  `ALTER_INSTRUCTION` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '变更说明',
  `CONFIRM_BY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '确认人',
  `CONFIRM_OPINION` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '确认意见',
  `CREATION_MODE` int(0) DEFAULT NULL COMMENT '创建方式(1:人工录入;2:系统产生;3:信息修改)',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产变更' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for asset_basic
-- ----------------------------
DROP TABLE IF EXISTS `asset_basic`;
CREATE TABLE `asset_basic`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `DEVICE_TYPE` int(0) NOT NULL COMMENT '设备类型',
  `MANAGEMENT_IP` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '管理IP',
  `DEVICE_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '设备名称',
  `SERIAL_NUMBER` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '序列号',
  `DEVICE_STATUS` int(0) DEFAULT NULL COMMENT '状态(0:全部,1:待上线,2:已上线,3:下线,4:报废)',
  `MANAGED_STATE` int(0) DEFAULT NULL COMMENT '纳管状态',
  `DEVICE_PRODUCER` int(0) NOT NULL COMMENT '厂商',
  `DEVICE_MODEL` int(0) NOT NULL COMMENT '型号',
  `SUBTYPE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '子类型',
  `OUTLINE_STRUCTURE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '外形结构',
  `SPECIFICATIONS` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '规格',
  `U_NUMBER` int(0) DEFAULT NULL COMMENT 'U数',
  `USE_STORAGE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '使用存储',
  `DATACENTER_ID` int(0) DEFAULT NULL COMMENT '数据中心',
  `RELATED_SERVICE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关联业务',
  `SERVICE_PATH` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '业务路径',
  `DEVICE_MANAGER_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '设备负责人1',
  `DEVICE_MANAGER_TWO` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '设备负责人2',
  `BUSINESS_MANAGER_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '业务负责人1',
  `BUSINESS_MANAGER_TWO` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '业务负责人2',
  `OPERATING_SYSTEM` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作系统',
  `REMARK` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `AFFILIATED_ORGANIZATION` int(0) DEFAULT NULL COMMENT '所属组织机构',
  `EQUIPMENT_ROOM` int(0) NOT NULL COMMENT '所在机房',
  `OWNING_CABINET` int(0) DEFAULT NULL COMMENT '所在机柜',
  `REGION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '所在区域',
  `CABINET_LOCATION` int(0) DEFAULT NULL COMMENT '机柜位置',
  `ABREAST` int(0) DEFAULT NULL COMMENT '并排放置(0:否,1:是)',
  `LOCATION_DESCRIPTION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '位置描述',
  `EXTENSION_TEST` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '扩展测试',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产详情' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for asset_expansion
-- ----------------------------
DROP TABLE IF EXISTS `asset_expansion`;
CREATE TABLE `asset_expansion`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSET_ID` bigint(0) NOT NULL COMMENT '资产ID',
  `CONFIG_CATEGORY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '配置类别(1:基本信息,2:硬件配置,3:网络配置)',
  `PROPERTY_CATEGORY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '属性类别',
  `GROUP_ID` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分组ID',
  `PROPERTY_NAME_CN` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '属性名称',
  `PROPERTY_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '英文名称',
  `PROPERTY_VALUE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '属性值',
  `ASSOCIATED_TABLE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关联表名',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产扩展' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for asset_maintenance
-- ----------------------------
DROP TABLE IF EXISTS `asset_maintenance`;
CREATE TABLE `asset_maintenance`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSET_ID` int(0) NOT NULL COMMENT '资产ID',
  `MAINTENANCE_TYPE` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '维保类型（数据字典）',
  `MAINTENANCE_PROVIDER` int(0) DEFAULT NULL COMMENT '维保商',
  `START_AT` int(0) DEFAULT NULL COMMENT '开始日期',
  `FINISH_AT` int(0) DEFAULT NULL COMMENT '结束日期',
  `MAINTENANCE_PERIOD` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '维保期限',
  `PRODUCTION_AT` int(0) DEFAULT NULL COMMENT '出厂日期',
  `VERSION` int(0) NOT NULL COMMENT '版本号',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产维保' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for asset_management
-- ----------------------------
DROP TABLE IF EXISTS `asset_management`;
CREATE TABLE `asset_management`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSET_ID` int(0) NOT NULL COMMENT '资产ID',
  `ASSET_CODE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产编号',
  `SHUTDOWN_LEVEL` int(0) DEFAULT NULL COMMENT '关机级别',
  `SERVICE_LEVEL` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务级别',
  `SERVICE_CODE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务代码',
  `BELONG_DEPT` int(0) DEFAULT NULL COMMENT '所属部门',
  `EQUIPMENT_USE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '设备用途',
  `USER_DEPARTMENT` int(0) DEFAULT NULL COMMENT '使用部门',
  `USING_SITE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '使用地点',
  `VERSION` int(0) NOT NULL COMMENT '版本号',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for asset_tree
-- ----------------------------
DROP TABLE IF EXISTS `asset_tree`;
CREATE TABLE `asset_tree`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `STATUS` int(0) NOT NULL COMMENT '资产状态',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `MANAGEMENT_IP` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '管理IP',
  `SERIAL_NUMBER` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '序列号',
  `PROPERTY_ID` int(0) DEFAULT NULL COMMENT '属性ID',
  `PARENT_ID` int(0) DEFAULT NULL COMMENT '父ID',
  `TYPE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型',
  `REMARK` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产树' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for assets_directory
-- ----------------------------
DROP TABLE IF EXISTS `assets_directory`;
CREATE TABLE `assets_directory`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `parent_id` bigint(0) NOT NULL COMMENT '父节点',
  `sort` bigint(0) DEFAULT NULL COMMENT '序号',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `created_at` bigint(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `updated_at` bigint(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `deleted_at` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产目录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for assets_expansion
-- ----------------------------
DROP TABLE IF EXISTS `assets_expansion`;
CREATE TABLE `assets_expansion`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `assets_id` bigint(0) NOT NULL COMMENT '资产id',
  `config_category` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '配置类别',
  `group_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分组ID',
  `name_cn` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '属性名称',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '英文名称',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '属性值',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `created_at` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `updated_at` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `deleted_at` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资产扩展-西航' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for cabinet_group
-- ----------------------------
DROP TABLE IF EXISTS `cabinet_group`;
CREATE TABLE `cabinet_group`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `CABINET_GROUP_CODE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '机柜组编号',
  `ROOM_ID` int(0) NOT NULL COMMENT '所属机房',
  `CABINET_GROUP_TYPE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '机柜组类型',
  `ROW` int(0) NOT NULL COMMENT '行',
  `START_COLUMN` int(0) NOT NULL COMMENT '开始列',
  `COLUMN` int(0) DEFAULT NULL COMMENT '所在列',
  `DUTY_PERSON_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '责任人1',
  `DUTY_PERSON_TWO` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人2',
  `USE_NOTES` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用途说明',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '机柜组信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for component_type
-- ----------------------------
DROP TABLE IF EXISTS `component_type`;
CREATE TABLE `component_type`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `COMPONENT_TYPE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '部件类型',
  `REMARK` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `COMPONENT_PICTURE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '部件图',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '部件类型' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for computer_room
-- ----------------------------
DROP TABLE IF EXISTS `computer_room`;
CREATE TABLE `computer_room`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ROOM_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `ROOM_CODE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '编码',
  `IDC_LOCATION` int(0) DEFAULT NULL COMMENT '所在IDC',
  `SUBGALLERY` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '所属楼座',
  `FLOOR` int(0) DEFAULT NULL COMMENT '所属楼层',
  `VOLTAGE` int(0) DEFAULT NULL COMMENT '电压',
  `ELECTRIC` int(0) DEFAULT NULL COMMENT '电流',
  `ROW_MAX` int(0) NOT NULL COMMENT '最大行数',
  `COLUMN_MAX` int(0) NOT NULL COMMENT '最大列数',
  `CABINET_NUMBER` int(0) NOT NULL COMMENT '可容纳机柜数',
  `ROOM_BEARING_CAPACITY` decimal(24, 6) DEFAULT NULL COMMENT '机房承重',
  `ROOM_AREA` decimal(24, 6) DEFAULT NULL COMMENT '机房面积',
  `RATED_POWER` int(0) DEFAULT NULL COMMENT '额定功率',
  `ROOM_PICTURE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '机房图片',
  `DUTY_PERSON_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人1',
  `DUTY_PERSON_two` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人2',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '机房信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dashboard_user
-- ----------------------------
DROP TABLE IF EXISTS `dashboard_user`;
CREATE TABLE `dashboard_user`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint(0) NOT NULL COMMENT '用户id',
  `assets_id` bigint(0) NOT NULL COMMENT '资产id',
  `type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '资产类型',
  `page_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '页签',
  `sort` int(0) DEFAULT NULL COMMENT '序号',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `created_at` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `updated_at` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `deleted_at` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '凤八大屏数据看板' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_cabinet
-- ----------------------------
DROP TABLE IF EXISTS `device_cabinet`;
CREATE TABLE `device_cabinet`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `CABINET_ID` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '机柜ID',
  `BELONG_ROOM` int(0) NOT NULL COMMENT '所在机房',
  `CABINET_CODE` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '机柜编号',
  `CABINET_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '机柜名称',
  `PRODUCER_ID` int(0) DEFAULT NULL COMMENT '厂商',
  `CABINET_MODEL` int(0) DEFAULT NULL COMMENT '型号',
  `CABINET_PICTURE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '机柜图片',
  `UNUMBER` int(0) NOT NULL COMMENT '规格(U数)',
  `ROW_NUMBER` int(0) DEFAULT NULL COMMENT '所在行',
  `ROW_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '所在行名称',
  `COLUMN_NUMBER` int(0) DEFAULT NULL COMMENT '所在列',
  `COLUMN_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '所在列名称',
  `BELONG_GROUP` int(0) DEFAULT NULL COMMENT '所属机柜组',
  `MAIN_POWER_SUPPLY` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '主要供电来源',
  `STANDBY_POWER_SUPPLY` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '临时供电来源',
  `POWER_SUPPLY_MODE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '供电方式',
  `POWER_CONSUMPTION` int(0) DEFAULT NULL COMMENT '电源功耗',
  `RATED_VOLTAGE` int(0) DEFAULT NULL COMMENT '额定电压',
  `RATED_CURRENT` int(0) DEFAULT NULL COMMENT '额定电流',
  `USAGE` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用途',
  `CABINET_TYPE` int(0) NOT NULL COMMENT '机柜类型;1:大一体机机柜;2:普通机柜;3:屏蔽机柜',
  `RESERVED_CABINET` int(0) NOT NULL COMMENT '预留机柜',
  `UNAVAILABLE_SPACE` int(0) NOT NULL COMMENT '不可用空间',
  `DUTY_PERSON_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人1',
  `DUTY_PERSON_TWO` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人2',
  `CABINET_BEARING_CAPACITY` decimal(24, 6) NOT NULL COMMENT '机柜承重',
  `CABINET_AREA` decimal(24, 6) NOT NULL COMMENT '机柜面积',
  `SERVICE_PARTITION` int(0) DEFAULT NULL COMMENT '业务分区',
  `POWER_PLUG_NUMBER` int(0) NOT NULL COMMENT '电源插头数量',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '机柜信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_model
-- ----------------------------
DROP TABLE IF EXISTS `device_model`;
CREATE TABLE `device_model`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `NAME` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '型号名称',
  `DEVICE_TYPE` int(0) DEFAULT NULL COMMENT '设备类型',
  `SUBTYPE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '子类型',
  `PRODUCER_ID` int(0) DEFAULT NULL COMMENT '厂商',
  `MODEL` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '型号',
  `SERIES` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '系列',
  `U_NUMBER` int(0) DEFAULT NULL COMMENT 'U数',
  `OUTLINE_STRUCTURE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '外形结构',
  `SPECIFICATIONS` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '规格',
  `MAXIMUM_MEMORY` decimal(24, 6) DEFAULT NULL COMMENT '最大内存(M)',
  `WORKING_CONSUMPTION` decimal(24, 6) DEFAULT NULL COMMENT '工作功耗(W)',
  `RATED_CONSUMPTION` decimal(24, 6) DEFAULT NULL COMMENT '额定功耗(W)',
  `PEAK_CONSUMPTION` decimal(24, 6) DEFAULT NULL COMMENT '峰值功耗(W)',
  `WEIGHT` decimal(24, 6) DEFAULT NULL COMMENT '设备重量(kg)',
  `ENLISTMENT` int(0) DEFAULT NULL COMMENT '服役期限(月)',
  `OUT_BAND_VERSION` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '带外版本',
  `DESCRIBE` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
  `PICTURE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上传照片',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备型号' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_online
-- ----------------------------
DROP TABLE IF EXISTS `device_online`;
CREATE TABLE `device_online`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `DEVICE_STATUS` int(0) NOT NULL COMMENT '类型',
  `ASSET_ID` int(0) NOT NULL COMMENT '资产ID',
  `DESCRIPTION` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '说明',
  `LINE_AT` int(0) DEFAULT NULL COMMENT '上线/下线日期',
  `LINE_DIRECTORY` int(0) DEFAULT NULL COMMENT '上线/下线目录',
  `CREATED_BY` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备上线下线记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_producer
-- ----------------------------
DROP TABLE IF EXISTS `device_producer`;
CREATE TABLE `device_producer`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `PRODUCER_TYPE` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '厂商类型',
  `ALIAS` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '简称',
  `CHINESE_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '中文名称',
  `COMPANY_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '公司全称',
  `SERVICE_TEL` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务电话',
  `SERVICE_EMAIL` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务邮箱',
  `COUNTRY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '国家',
  `CITY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '城市',
  `ADDRESS` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '地址',
  `FAX` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '传真',
  `CONTACT_PERSON` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '联系人',
  `CONTACT_NUMBER` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '联系人电话',
  `CONTACT_EMAIL` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '联系人邮箱',
  `OFFICIAL` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '官方站点',
  `IS_DOMESTIC` int(0) DEFAULT NULL COMMENT '是否国产',
  `IS_DISPLAY_CHINESE` int(0) DEFAULT NULL COMMENT '是否显示中文',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备厂商' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_scrap
-- ----------------------------
DROP TABLE IF EXISTS `device_scrap`;
CREATE TABLE `device_scrap`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSET_ID` int(0) NOT NULL COMMENT '资产ID',
  `DEVICE_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '设备名称',
  `SERIAL_NUMBER` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '序列号',
  `OLD_MANAGEMENT_IP` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '管理IP',
  `DEVICE_PRODUCER` int(0) NOT NULL COMMENT '厂商',
  `DEVICE_TYPE` int(0) NOT NULL COMMENT '设备类型',
  `DEVICE_MODEL` int(0) NOT NULL COMMENT '型号',
  `ASSET_CODE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产编号',
  `OLD_DATACENTER` int(0) DEFAULT NULL COMMENT '原数据中心',
  `OLD_LOCATION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原所在位置',
  `PURCHASE_AT` int(0) DEFAULT NULL COMMENT '采购日期',
  `OLD_DEVICE_MANAGER` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原责任人',
  `OLD_BELONG_ORGANIZATION` int(0) DEFAULT NULL COMMENT '所属组织机构',
  `REMARK` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '报废说明',
  `SCRAP_AT` int(0) DEFAULT NULL COMMENT '报废时间',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备报废' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_type
-- ----------------------------
DROP TABLE IF EXISTS `device_type`;
CREATE TABLE `device_type`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `TYPES` int(0) NOT NULL COMMENT '类别(1:设备类型;2:备件设备类型)',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备类型' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device_type_config
-- ----------------------------
DROP TABLE IF EXISTS `device_type_config`;
CREATE TABLE `device_type_config`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '名称',
  `TYPE` int(0) DEFAULT NULL COMMENT '设备类型',
  `TYPE_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '设备类型名称',
  `CONFIG` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '表单属性配置',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备类型表单配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dict_data
-- ----------------------------
DROP TABLE IF EXISTS `dict_data`;
CREATE TABLE `dict_data`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `TYPE_CODE` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典编码',
  `DICT_KEY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典标签',
  `DICT_VALUE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典键值',
  `SN` int(0) DEFAULT NULL COMMENT '序号',
  `REMARK` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典数据表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dict_type
-- ----------------------------
DROP TABLE IF EXISTS `dict_type`;
CREATE TABLE `dict_type`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `TYPE_CODE` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典编码',
  `DICT_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典名称',
  `IS_VISIBLE` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '是否可见',
  `REMARK` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` bigint(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典类别表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for distribution_frame
-- ----------------------------
DROP TABLE IF EXISTS `distribution_frame`;
CREATE TABLE `distribution_frame`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ROOM_ID` int(0) NOT NULL COMMENT '所属机房',
  `CABINET_ID` int(0) NOT NULL COMMENT '所属机柜',
  `DIS_FRAME_CODE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '配线架编号',
  `DIS_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '配线架名称',
  `PRODUCER_ID` int(0) DEFAULT NULL COMMENT '厂商',
  `MODEL` int(0) DEFAULT NULL COMMENT '型号',
  `SPECIFICATION` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '规格',
  `DIS_TYPE` int(0) DEFAULT NULL COMMENT '配线架类型(0:双绞线;1:光纤配线架)',
  `TOTAL_PORT_NUM` int(0) NOT NULL COMMENT '总端口数',
  `USED_PORT_NUM` int(0) DEFAULT NULL COMMENT '已用端口数',
  `PORT_PREFIX` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '端口前缀',
  `CABINET_LOCATION` int(0) DEFAULT NULL COMMENT '机柜位置(U)',
  `PURCHASE_AT` int(0) DEFAULT NULL COMMENT '采购日期',
  `DIS_PICTURE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '配线架图片',
  `USE` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用途',
  `DUTY_PERSON_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '责任人1',
  `DUTY_PERSON_TWO` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人2',
  `UNUMBER` int(0) DEFAULT NULL COMMENT 'U数',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '配线架信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for feorg
-- ----------------------------
DROP TABLE IF EXISTS `feorg`;
CREATE TABLE `feorg`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `parentid` int(0) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `parent_id` int(0) DEFAULT NULL,
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for maintenance_service_config
-- ----------------------------
DROP TABLE IF EXISTS `maintenance_service_config`;
CREATE TABLE `maintenance_service_config`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `MAINTENANCE_ID` int(0) NOT NULL COMMENT '维保ID',
  `SERVICE_OPTION_CODE` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务选项编码',
  `SERVICE_OBJ_KEY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务对象标签',
  `SERVICE_OBJ_VALUE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '服务对象值',
  `DEADLINE` int(0) DEFAULT NULL COMMENT '服务截止时间',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '维保服务项配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for monitoring
-- ----------------------------
DROP TABLE IF EXISTS `monitoring`;
CREATE TABLE `monitoring`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSET_ID` bigint(0) DEFAULT NULL COMMENT '资产id',
  `MONITORING_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '监控名称',
  `DATASOURCE_ID` bigint(0) NOT NULL COMMENT '数据源id',
  `MONITORING_SQL` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '监控脚本',
  `STATUS` int(0) NOT NULL COMMENT '状态',
  `IS_ALARM` int(0) NOT NULL COMMENT '是否启用告警(0:未启动；1:启动)',
  `TARGET_ID` int(0) NOT NULL COMMENT '采集器',
  `CONFIG` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '配置信息',
  `REMARK` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '说明',
  `UNIT` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '单位',
  `LABEL` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '标签',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '监控' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '日志主键',
  `TYPE` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型',
  `OBJECT` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '0' COMMENT '对象',
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
  `USER` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户',
  `OPER_TIME` int(0) DEFAULT 0 COMMENT '执行时间',
  `OPER_URL` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求URL',
  `OPER_PARAM` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求参数',
  `JSON_RESULT` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '返回参数',
  `REQ_METHOD` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求方式',
  `STATUS` int(0) DEFAULT 0 COMMENT '操作状态码',
  `ERROR_MSG` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '错误消息',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '操作日志' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for pdu
-- ----------------------------
DROP TABLE IF EXISTS `pdu`;
CREATE TABLE `pdu`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ASSETS_CODE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '资产编号',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '名称',
  `BRAND` int(0) DEFAULT NULL COMMENT '品牌',
  `MODEL` int(0) DEFAULT NULL COMMENT '型号',
  `STANDARD` int(0) NOT NULL COMMENT '标准(1:新国标;2:国标)',
  `JACK_NUM` int(0) NOT NULL COMMENT '插孔数',
  `LIMIT_VOLTAGE` decimal(24, 6) DEFAULT NULL COMMENT '限制电压(V)',
  `MAX_ELECTRIC` decimal(24, 6) DEFAULT NULL COMMENT '最大耐冲击电压(KA)',
  `USE` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用途',
  `PURCHASE_AT` int(0) DEFAULT NULL COMMENT '采购日期',
  `POWER` decimal(24, 6) DEFAULT NULL COMMENT '功率',
  `UNIT_PRICE` decimal(24, 6) DEFAULT NULL COMMENT '单价',
  `BELONG_ROOM` int(0) NOT NULL COMMENT '所在机房',
  `CABINET_ID` int(0) NOT NULL COMMENT '所在机柜编号',
  `DUTY_PERSON_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '责任人1',
  `DUTY_PERSON_TWO` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人2',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'PDU' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for room_partition
-- ----------------------------
DROP TABLE IF EXISTS `room_partition`;
CREATE TABLE `room_partition`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ROOM_ID` int(0) NOT NULL COMMENT '机房ID',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分区名称',
  `START_ROW` int(0) NOT NULL COMMENT '起始行',
  `START_COLUMN` int(0) NOT NULL COMMENT '起始列',
  `HEIGHT` int(0) NOT NULL COMMENT '高度',
  `WIDTH` int(0) NOT NULL COMMENT '宽度',
  `SPACE_TYPE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '空间类型',
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '位置描述',
  `DUTY_PERSON_ONE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人1',
  `DUTY_PERSON_two` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人2',
  `REMARK` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '机房分区表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for spare_part_basic
-- ----------------------------
DROP TABLE IF EXISTS `spare_part_basic`;
CREATE TABLE `spare_part_basic`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `PRODUCT_ID` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '商品编号',
  `COMPONENT_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '部件名称',
  `COMPONENT_TYPE` int(0) NOT NULL COMMENT '部件类型',
  `COMPONENT_NUM` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部件号',
  `COMPONENT_BRAND` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '部件品牌',
  `SPECIFICATION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '型号规格',
  `COMPONENT_UNIT` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部件单位',
  `UNIT_PRICE` decimal(24, 6) DEFAULT NULL COMMENT '单价(元)',
  `DEVICE_TYPE` int(0) NOT NULL COMMENT '设备类型',
  `ASSET_CLASSIFICATION` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资产分类',
  `BELONG_LINE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '所属条线',
  `BELONG_ORGANIZATION` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '所属单位',
  `PURCHASING_APPLICANT` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '采购申请人',
  `SUPPLIER` int(0) DEFAULT NULL COMMENT '供应商',
  `DETAILED_CONFIG` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '详细配置',
  `REMARK` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `SPARE_PART_DETAIL` int(0) NOT NULL COMMENT '备件明细(0:否;1:是)',
  `COMPONENT_PICTURE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部件图片',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '备件基础数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for storeroom_management
-- ----------------------------
DROP TABLE IF EXISTS `storeroom_management`;
CREATE TABLE `storeroom_management`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ROOM_NUMBER` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '房间号',
  `BELONG_IDC` int(0) NOT NULL COMMENT '所属IDC',
  `ROOM_ADDRESS` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '房间地址',
  `DUTY_BY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '责任人',
  `SHELF_INFORMATION` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '货架信息',
  `CONTACT_NUMBER` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '联系电话',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '库房信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_config
-- ----------------------------
DROP TABLE IF EXISTS `user_config`;
CREATE TABLE `user_config`  (
  `ID` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '配置主键',
  `LOG_LEVER` int(0) DEFAULT NULL COMMENT '日志等级',
  `HTTP_HOST` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'HTTP监听地址',
  `HTTP_PORT` int(0) NOT NULL DEFAULT 0 COMMENT 'HTTP监听端口',
  `CAPTCHA` int(0) NOT NULL DEFAULT 0 COMMENT '启用验证码',
  `API_SERVICE` int(0) NOT NULL DEFAULT 0 COMMENT 'APIForService',
  `ACCESS_EXPIRED` int(0) NOT NULL DEFAULT 0 COMMENT 'token有效期（用户登录）',
  `REFRESH_EXPIRED` int(0) NOT NULL DEFAULT 0 COMMENT 'token有效期（长期展示）',
  `OPEN_RSA` int(0) NOT NULL DEFAULT 0 COMMENT '启用加密',
  `LOGIN_TITLE` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '登录页标题',
  `LOGO_TOP` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '系统顶部LOGO',
  `LOGO_TITLE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '网页标题LOGO',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_config
-- ----------------------------
INSERT INTO `user_config` VALUES (1, 1, '0.0.0.0', 17000, 2, 1, 999999, 10080, 2, 'yugu', 'images/logo_top.png', 'images/logo_title.png', '0', 0, '0', 1700200695, NULL);


-- ----------------------------
-- Table structure for api_service
-- ----------------------------
DROP TABLE IF EXISTS `api_service`;
CREATE TABLE `api_service`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` int(0) DEFAULT NULL COMMENT '删除时间',
  `NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `TYPE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型;sql or promql',
  `DATASOURCE_ID` int(0) DEFAULT NULL COMMENT '数据源;promql 需要指定数据源',
  `URL` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'URL',
  `SCRIPT` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '执行脚本',
  `VALUE_FIELD` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '值字段',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '接口管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of api_service
-- ----------------------------
INSERT INTO `api_service` VALUES (1, 'root', 1700530855, 'root', 1700550729, 0, '资产品牌分布', 'sql', 0, '/asset/brand', 'select count(*) value , type name from assets group by type', '');
INSERT INTO `api_service` VALUES (2, 'root', 1701165238, 'root', 1701165277, 0, '资产健康度', 'sql', 0, '/asset_health', 'select t.alert name,count(*) value from (\nselect a.id id, if(exists (select 1 from alert_cur_event ace where ace.asset_id = a.id and ace.deleted_at is null),\'告警\',\'正常\') alert\nfrom  assets a where  a.deleted_at is null ) t\ngroup by t.alert', '');
INSERT INTO `api_service` VALUES (3, 'root', 1701165847, 'root', 1701165881, 0, '资产监控状态', 'sql', 0, '/asset_monitor', 'select t.status name,count(*) value from (\nselect a.id id, if(a.`status` = 1,\'正常\',\'离线\') status\nfrom  assets a where  a.deleted_at is null ) t\ngroup by t.status', '');
INSERT INTO `api_service` VALUES (4, 'root', 1701221592, 'root', 1701221605, 0, '资产告警状态', 'sql', 0, '/asset_alert', 'select t.alert name, t.type, count(*) value from (\nselect a.id id, a.type type, if(exists (select 1 from alert_cur_event ace where ace.asset_id = a.id and ace.deleted_at is null),\'告警\',\'正常\') alert\nfrom  assets a where  a.deleted_at is null ) t\ngroup by t.alert, t.type', '');
INSERT INTO `api_service` VALUES (5, 'root', 1701225213, 'root', 1701225219, 0, '7日内告警统计', 'sql', 0, '/alerts_in_week', 'select dat name, count(*) value from (\nselect FROM_UNIXTIME(trigger_time, \'%Y-%m-%d\') dat from alert_his_event where trigger_time > UNIX_TIMESTAMP(date_sub(now(), interval 7 day)) )t group by t.dat', '');
INSERT INTO `api_service` VALUES (6, 'root', 1701226840, '', 1701226840, 0, '品牌故障率', 'sql', 0, '/asset_fault', '\nselect t.brand name, sum(t.error)/count(*) value  from\n(select id, if(manufacturers is null or manufacturers = \'\', \'其他\', manufacturers) brand ,\nif( exists (select 1 from alert_cur_event ace  where ace.asset_id = a.ID), 1 , 0) error\nfrom assets a\nwhere deleted_at is null ) t\ngroup by t.brand', '');


-- ----------------------------
-- Table structure for bigscreen
-- ----------------------------
DROP TABLE IF EXISTS `bigscreen`;
CREATE TABLE `bigscreen`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  `TITLE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
  `DESC` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '简介',
  `CONFIG` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '配置',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bigscreen
-- ----------------------------
INSERT INTO `bigscreen` VALUES (1, 'root', 1700639591, 'root', 1701413642, NULL, '首页大屏', '首页大屏', '{\"title\":\"首页大屏\",\"gridSize\":10,\"gridBorderColor\":\"#090548\",\"showAuxiliary\":true,\"horizontalNumber\":3,\"verticalNumber\":3,\"interval\":10,\"auxiliaryBorderColor\":\"#1890ff\",\"width\":1920,\"height\":1080,\"backgroundImage\":\"\",\"backgroundColor\":\"#090548\",\"widgets\":[{\"id\":\"fb966a94-cc4a-48a5-bffa-b5a46dc32a3c\",\"code\":\"BorderBox11\",\"configureValue\":{\"title\":\"工控网运维监控系统\",\"titleWidth\":250},\"coordinateValue\":{\"width\":1921,\"height\":1080,\"x\":0,\"y\":0}},{\"id\":\"0e58992d-fac3-420e-9428-07056dacc966\",\"code\":\"Line\",\"configureValue\":{\"styleAnimateInfinite\":false,\"styleAnimationDelay\":0,\"styleAnimationName\":\"\",\"styleAnimationDuration\":1,\"styleAnimationTimingFunction\":\"linear\",\"titleTextShow\":false,\"titleText\":\"\",\"titleTextFontSize\":14,\"titleTextLineHeight\":1.2,\"titleTextFontFamily\":\"Microsoft YaHei\",\"titleTextFontWeight\":\"bold\",\"titleTextColor\":\"#fff\",\"legendShow\":true,\"legendType\":\"plain\",\"legendOrient\":\"horizontal\",\"legendFontSize\":12,\"legendIcon\":\"rect\",\"legendColor\":\"#fff\",\"legendLeft\":\"center\",\"legendTop\":\"top\",\"gridShow\":false,\"gridLeft\":50,\"gridRight\":30,\"gridTop\":30,\"gridBottom\":30,\"gridBorderColor\":\"#ccc\",\"xAxisShow\":true,\"xAxisType\":\"category\",\"xAxisName\":\"\",\"xAxisNameLocation\":\"end\",\"xAxisNameTextStyleFontSize\":12,\"xAxisNameTextStyleLineHeight\":12,\"xAxisNameTextStyleFontFamily\":\"serif\",\"xAxisNameTextStyleFontWeight\":\"normal\",\"xAxisBoundaryGap\":false,\"xAxisNameRotate\":0,\"xAxisLineShow\":true,\"xAxisLabelShow\":true,\"xAxisLabelRotate\":0,\"xAxisSplitLineShow\":true,\"xAxisSplitAreaShow\":false,\"xAxisSplitAreaOpacity\":10,\"xAxisPointerShow\":true,\"xAxisTickShow\":true,\"xAxisAlignWithLabel\":false,\"yAxisShow\":true,\"yAxisType\":\"value\",\"yAxisName\":\"\",\"yAxisNameLocation\":\"end\",\"yAxisNameTextStyleFontSize\":12,\"yAxisNameTextStyleLineHeight\":12,\"yAxisNameTextStyleFontFamily\":\"serif\",\"yAxisNameTextStyleFontWeight\":\"normal\",\"yAxisBoundaryGap\":false,\"yAxisNameRotate\":0,\"yAxisLineShow\":true,\"yAxisLabelShow\":true,\"yAxisLabelRotate\":0,\"yAxisSplitLineShow\":true,\"yAxisSplitAreaShow\":false,\"yAxisSplitAreaOpacity\":10,\"yAxisPointerShow\":false,\"yAxisTickShow\":true,\"yAxisAlignWithLabel\":false,\"axisNameColor\":\"rgba(255,255,255,.2)\",\"axisLineColor\":\"rgba(255,255,255,.2)\",\"axisLabelColor\":\"rgba(255,255,255,.8)\",\"splitLineColor\":\"rgba(255,255,255,.2)\",\"axisPointerColor\":\"red\",\"themeColor1\":\"#fc97af\",\"themeColor2\":\"#87f7cf\",\"themeColor3\":\"#f7f494\",\"themeColor4\":\"#72ccff\",\"themeColor5\":\"#f7c5a0\",\"themeColor6\":\"#d4a4eb\",\"themeColor7\":\"#d2f5a6\",\"themeColor8\":\"#76f2f2\",\"lineWidth\":2,\"lineSmooth\":false,\"lineAreaStyle\":false,\"lineAreaStyleOpacity\":70,\"seriesLabelShow\":false,\"seriesLabelPosition\":\"top\",\"seriesLabelColor\":\"#fff\",\"seriesStackValue\":\"\",\"showSymbol\":true,\"symbol\":\"circle\",\"symbolSize\":4},\"coordinateValue\":{\"width\":500,\"height\":300,\"x\":44,\"y\":109},\"dataValue\":{\"dataType\":\"api\",\"field\":\"series\",\"autoRefresh\":true,\"url\":\"/api/n9e/api-service/2/execute\",\"interval\":15}},{\"id\":\"fc67bd59-ace8-4615-8ca4-6e7522ee0bcd\",\"code\":\"Pie\",\"configureValue\":{\"styleDisplay\":\"block\",\"styleAnimateInfinite\":false,\"styleAnimationDelay\":0,\"styleAnimationName\":\"\",\"styleAnimationDuration\":1,\"styleAnimationTimingFunction\":\"linear\",\"titleTextShow\":false,\"titleText\":\"\",\"titleTextFontSize\":14,\"titleTextLineHeight\":1.2,\"titleTextFontFamily\":\"Microsoft YaHei\",\"titleTextFontWeight\":\"bold\",\"titleTextColor\":\"#fff\",\"legendShow\":true,\"legendType\":\"plain\",\"legendOrient\":\"horizontal\",\"legendFontSize\":12,\"legendIcon\":\"rect\",\"legendColor\":\"#fff\",\"legendLeft\":\"center\",\"legendTop\":\"top\",\"axisNameColor\":\"rgba(255,255,255,.2)\",\"axisLineColor\":\"rgba(255,255,255,.2)\",\"axisLabelColor\":\"rgba(255,255,255,.8)\",\"splitLineColor\":\"rgba(255,255,255,.2)\",\"axisPointerColor\":\"red\",\"themeColor1\":\"#fc97af\",\"themeColor2\":\"#87f7cf\",\"themeColor3\":\"#f7f494\",\"themeColor4\":\"#72ccff\",\"themeColor5\":\"#f7c5a0\",\"themeColor6\":\"#d4a4eb\",\"themeColor7\":\"#d2f5a6\",\"themeColor8\":\"#76f2f2\",\"seriesInsideRadius\":0,\"seriesAutsideRadius\":80,\"seriesRoseType\":false,\"xAxisShow\":false,\"yAxisShow\":false,\"seriesLabelShow\":true,\"seriesLabelPosition\":\"outside\",\"seriesLabelColor\":\"\",\"xAxisType\":\"category\",\"yAxisType\":\"value\"},\"coordinateValue\":{\"width\":500,\"height\":300,\"x\":857,\"y\":506},\"dataValue\":{\"dataType\":\"api\",\"field\":\"series\",\"autoRefresh\":true,\"url\":\"/api/n9e/api-service/3/execute\",\"interval\":15}},{\"id\":\"464a9799-a7e3-46ed-870b-5d729b09d2a7\",\"code\":\"Decoration1\",\"configureValue\":{},\"coordinateValue\":{\"width\":200,\"height\":50,\"x\":19,\"y\":44}},{\"id\":\"a439b611-a96c-4155-96cd-e1df8dd250b7\",\"code\":\"BaseText\",\"configureValue\":{\"styleFontSize\":32,\"styleLetterSpacing\":0,\"styleFontWeight\":\"bold\",\"styleTextAlign\":\"center\",\"styleBackgroundColor\":\"\",\"styleFontFamily\":\"Microsoft YaHei\",\"styleLineHeight\":1,\"styleColor\":\"rgba(239,226,68,1)\",\"styleBoxInset\":false,\"styleBoxShadowX\":0,\"styleBoxShadowY\":0,\"styleBoxShadowF\":0,\"styleBoxShadowC\":\"\",\"styleBorderStyle\":\"none\",\"styleBorderWidth\":0,\"styleBorderColor\":\"\",\"styleBorderTopLeftRadius\":0,\"styleBorderTopRightRadius\":0,\"styleBorderBottomLeftRadius\":0,\"styleBorderBottomRightRadius\":0,\"styleAnimateInfinite\":false,\"styleAnimationDelay\":0,\"styleAnimationName\":\"\",\"styleAnimationDuration\":1,\"styleAnimationTimingFunction\":\"linear\"},\"coordinateValue\":{\"width\":181,\"height\":46,\"x\":177,\"y\":46},\"dataValue\":{\"dataType\":\"mock\",\"field\":\"value\",\"autoRefresh\":false,\"mock\":{\"series\":[{\"seriesName\":\"Email\",\"data\":[{\"name\":\"Mon\",\"value\":\"@integer(100, 300)\"},{\"name\":\"Tue\",\"value\":\"@integer(100, 300)\"},{\"name\":\"Wed\",\"value\":\"@integer(100, 300)\"},{\"name\":\"Thu\",\"value\":\"@integer(100, 300)\"},{\"name\":\"Fri\",\"value\":\"@integer(100, 300)\"},{\"name\":\"Sat\",\"value\":\"@integer(100, 300)\"},{\"name\":\"Sun\",\"value\":\"@integer(100, 300)\"}]}],\"value\":\"纳管资产\"}}}],\"description\":\"这是一个测试大屏\"}');

-- ----------------------------
-- Table structure for license_config
-- ----------------------------
DROP TABLE IF EXISTS `license_config`;
CREATE TABLE `license_config`  (
  `ID` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `DAYS` int(0) DEFAULT NULL COMMENT '剩余天数',
  `NODES` int(0) DEFAULT NULL COMMENT '剩余节点数',
  `FREQUENCY` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '提醒频率',
  `EMAIL` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
  `CLIENT` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '客户',
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '许可配置' ROW_FORMAT = Dynamic;

INSERT INTO `license_config` VALUES (1, 10, 30, 'once', '', '', 'root', 1701243521, 'root', 1701312136, NULL);

CREATE TABLE `es_index_pattern` (
    `id` bigint unsigned not null auto_increment,
    `datasource_id` bigint not null default 0 comment 'datasource id',
    `name` varchar(191) not null,
    `time_field` varchar(128) not null default '@timestamp',
    `allow_hide_system_indices` tinyint(1) not null default 0,
    `fields_format` varchar(4096) not null default '',
    `create_at` bigint default '0',
    `create_by` varchar(64) default '',
    `update_at` bigint default '0',
    `update_by` varchar(64) default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`datasource_id`, `name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
