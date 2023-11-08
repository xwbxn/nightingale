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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

insert into user_group(id, name, create_at, create_by, update_at, update_by) values(1, 'demo-root-group', unix_timestamp(now()), 'root', unix_timestamp(now()), 'root');

CREATE TABLE `user_group_member` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint unsigned not null,
    `user_id` bigint unsigned not null,
    KEY (`group_id`),
    KEY (`user_id`),
    PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

insert into user_group_member(group_id, user_id) values(1, 1);

CREATE TABLE `configs` (
    `id` bigint unsigned not null auto_increment,
    `ckey` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `cval` text not null,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`ckey`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `role` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

insert into busi_group(id, name, create_at, create_by, update_at, update_by) values(1, 'Default Busi Group', unix_timestamp(now()), 'root', unix_timestamp(now()), 'root');

CREATE TABLE `busi_group_member` (
    `id` bigint unsigned not null auto_increment,
    `busi_group_id` bigint not null comment 'busi group id',
    `user_group_id` bigint not null comment 'user group id',
    `perm_flag` char(2) not null comment 'ro | rw',
    PRIMARY KEY (`id`),
    KEY (`busi_group_id`),
    KEY (`user_group_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

-- for dashboard new version
CREATE TABLE `board_payload` (
    `id` bigint unsigned not null comment 'dashboard id',
    `payload` mediumtext not null,
    UNIQUE KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

-- deprecated
-- auto create the first subclass 'Default chart group' of dashboard
CREATE TABLE `chart_group` (
    `id` bigint unsigned not null auto_increment,
    `dashboard_id` bigint unsigned not null,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `weight` int not null default 0,
    PRIMARY KEY (`id`),
    KEY (`dashboard_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

-- deprecated
CREATE TABLE `chart` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint unsigned not null comment 'chart group id',
    `configs` text,
    `weight` int not null default 0,
    PRIMARY KEY (`id`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `chart_share` (
    `id` bigint unsigned not null auto_increment,
    `cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `datasource_id` bigint unsigned not null default 0,
    `configs` text,
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    primary key (`id`),
    key (`create_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `alert_rule` (
    `id` bigint unsigned not null auto_increment,
    `group_id` bigint not null default 0 comment 'busi group id',
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
    `rule_config` text not null comment 'rule_config',
    `prom_ql` text not null comment 'promql',
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
    `annotations` text not null comment 'annotations',
    `extra_config` text not null comment 'extra_config',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`group_id`),
    KEY (`update_at`)
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;



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
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

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
    `query_configs` text not null comment 'query configs',
    `create_at` bigint default '0',
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '',
    `update_at` bigint default '0',
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default '',
    PRIMARY KEY (`id`),
    KEY `group_id` (`group_id`),
    KEY `update_at` (`update_at`)
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

insert into alert_aggr_view(name, rule, cate) values('By BusiGroup, Severity', 'field:group_name::field:severity', 0);
insert into alert_aggr_view(name, rule, cate) values('By RuleName', 'field:rule_name', 0);

CREATE TABLE `alert_cur_event` (
    `id` bigint unsigned not null comment 'use alert_his_event.id',
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
    `annotations` text not null comment 'annotations',
    `rule_config` text not null comment 'annotations',
    `tags` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'merge data_tags rule_tags, split by ,,',
    `status` tinyint(1) not null default 0 comment '状态',
    PRIMARY KEY (`id`),
    KEY (`hash`),
    KEY (`rule_id`),
    KEY (`trigger_time`, `group_id`),
    KEY (`notify_repeat_next`)
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `alert_his_event` (
    `id` bigint unsigned not null AUTO_INCREMENT,
    `is_recovered` tinyint(1) not null,
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
    `annotations` text not null comment 'annotations',
    `rule_config` text not null comment 'annotations',
    `status` tinyint(1) not null default 0 comment '状态',
    `handle_by` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `handle_at` bigint not null default 0,
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`hash`),
    KEY (`rule_id`),
    KEY (`trigger_time`, `group_id`)
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4;

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
    `script`    text         not null,
    `args`      varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `tags`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'split by space',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `update_at` bigint not null default 0,
    `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`group_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `task_tpl_host`
(
    `ii`   int unsigned NOT NULL AUTO_INCREMENT,
    `id`   int unsigned not null comment 'task tpl id',
    `host` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  not null comment 'ip or hostname',
    PRIMARY KEY (`ii`),
    KEY (`id`, `host`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
    `script`    text            not null,
    `args`      varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci    not null default '',
    `create_at` bigint not null default 0,
    `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    PRIMARY KEY (`id`),
    KEY (`create_at`, `group_id`),
    KEY (`create_by`),
    KEY (`event_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `alerting_engines`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `instance` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'instance identification, e.g. 10.9.0.9:9090',
    `datasource_id` bigint not null default 0 comment 'datasource id',
    `engine_cluster` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '' comment 'n9e-alert cluster',
    `clock` bigint not null,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
    `settings` text not null,
    `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `http` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `auth` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `created_at` bigint not null default 0,
    `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    `updated_at` bigint not null default 0,
    `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null default '',
    UNIQUE KEY (`name`),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `builtin_cate` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) not null,
    `user_id` bigint not null default 0,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `notify_tpl` (
    `id` bigint unsigned not null auto_increment,
    `channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `content` text not null,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`channel`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `sso_config` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci not null,
    `content` text not null,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- CREATE TABLE `organization` (
--   `id` int(10) NOT NULL AUTO_INCREMENT,
--   `name` varchar(50) DEFAULT NULL,
--   `parent_id` int(10) DEFAULT NULL,
--   `path` varchar(50) DEFAULT NULL,
--   `son` int DEFAULT NULL,
--   `city` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   `manger` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   `create_at` bigint NOT NULL DEFAULT '0',
--   `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
--   `update_at` bigint NOT NULL DEFAULT '0',
--   `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
--   `deleted_at` datetime DEFAULT NULL,
--   PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
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
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO `datasource`(`id`, `name`, `description`, `category`, `plugin_id`, `plugin_type`, `plugin_type_name`, `cluster_name`, `settings`, `status`, `http`, `auth`) VALUES (1, 'default', '', '', 0, 'prometheus', '', 'default', '{\"write_addr\":\"http://victoria-metrics:8428/api/v1/write\"}', 'enabled', '{\"timeout\":10000,\"dial_timeout\":0,\"tls\":{\"skip_tls_verify\":false},\"max_idle_conns_per_host\":0,\"url\":\"http://victoria-metrics:8428\",\"headers\":{}}', '{\"basic_auth\":false,\"basic_auth_user\":\"vm\",\"basic_auth_password\":\"vmdctbcab\"}');
