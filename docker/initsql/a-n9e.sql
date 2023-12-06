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
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`username`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- insert into `users`(id, username, nickname, password, roles, create_at, create_by, update_at, update_by) values(1, 'root', '超管', 'root.2020', 'Admin', unix_timestamp(now()), 'system', unix_timestamp(now()), 'system');
INSERT INTO `users` VALUES (1, 'root', '超管', '042c05fffc2f49ca29a76223f3a41e83', '', '', '', 'Admin', 1, 1, '{}', 0, 1698905269, 'system', 1698973348, 'root', NULL);


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

insert into user_group(id, name, create_at, create_by, update_at, update_by) values(1, 'demo-root-group', unix_timestamp(now()), 'root', unix_timestamp(now()), 'root');

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

insert into busi_group(id, name, create_at, create_by, update_at, update_by) values(1, 'Default Busi Group', unix_timestamp(now()), 'root', unix_timestamp(now()), 'root');

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

-- for dashboard new version
CREATE TABLE `board_payload` (
    `id` bigint unsigned not null comment 'dashboard id',
    `payload` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    UNIQUE KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

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
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
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
    `callbacks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: http://a.com/api/x http://a.com/api/y',
    `runbook_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `append_tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: service=n9e mod=api',
    `annotations` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'annotations',
    `extra_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'extra_config',
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
    `periodic_mutes` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `severities` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
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
    `severities` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `tags` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'json,map,tagkey->regexp|value',
    `redefine_severity` tinyint(1) default 0 comment 'is redefine severity?',
    `new_severity` tinyint(1) not null comment '0:Emergency 1:Warning 2:Notice',
    `redefine_channels` tinyint(1) default 0 comment 'is redefine channels?',
    `new_channels` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space: sms voice email dingtalk wecom',
    `user_group_ids` varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'split by space 1 34 5, notify cc to user_group_ids',
    `webhooks` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `extra_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null comment 'extra_config',
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

INSERT INTO `datasource`(`id`, `name`, `description`, `category`, `plugin_id`, `plugin_type`, `plugin_type_name`, `cluster_name`, `settings`, `status`, `http`, `auth`) VALUES (1, 'default', '', '', 0, 'prometheus', '', 'default', '{\"write_addr\":\"http://victoria-metrics:8428/api/v1/write\"}', 'enabled', '{\"timeout\":10000,\"dial_timeout\":0,\"tls\":{\"skip_tls_verify\":false},\"max_idle_conns_per_host\":0,\"url\":\"http://victoria-metrics:8428\",\"headers\":{}}', '{\"basic_auth\":false,\"basic_auth_user\":\"vm\",\"basic_auth_password\":\"vmdctbcab\"}');


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
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '属性值',
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
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '接口管理' ROW_FORMAT = Dynamic;

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
  `DELETED_AT` int(0) DEFAULT NULL COMMENT '删除时间',
  `TITLE` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
  `DESC` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '简介',
  `CONFIG` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '配置',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;


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
  `CREATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '创建人',
  `CREATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `UPDATED_BY` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '更新人',
  `UPDATED_AT` int(0) NOT NULL DEFAULT 0 COMMENT '更新时间',
  `DELETED_AT` datetime(0) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '许可配置' ROW_FORMAT = Dynamic;

INSERT INTO `license_config` VALUES (1, 10, 30, 'once', '', 'root', 1701243521, 'root', 1701312136, NULL);