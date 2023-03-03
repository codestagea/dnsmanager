create table dns_view (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) not null default 'any',
    `remark` varchar(100) not null default '',
    PRIMARY KEY (`id`),
    unique key `name` (`name`)
);

create table dns_zone (
    `id`          bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
    `zone`        varchar(100) NOT NULL DEFAULT '' COMMENT 'zone',
    `refresh`     int(11) NOT NULL DEFAULT '28800',
    `retry`       int(11) NOT NULL DEFAULT '14400',
    `expire`      int(11) NOT NULL DEFAULT '86400',
    `minimum`     int(11) NOT NULL DEFAULT '86400',
    `serial`      bigint(20) NOT NULL DEFAULT '1',
    `host_master` varchar(64) NOT NULL default '',
    `primary_ns`  varchar(64) NOT NULL DEFAULT 'ns.ddns.net.',
    `state`       varchar(20) NOT NULL DEFAULT 'running' COMMENT '状态',
    `remark`      varchar(100) NOT NULL DEFAULT '' COMMENT '描述',
    primary key (`id`),
    unique key (`zone`)
);

CREATE TABLE `dns_record`(
    `id`      bigint(20) unsigned NOT NULL AUTO_INCREMENT comment 'id',
    `zone_id` bigint(20) NOT NULL default 0,
    `host`    varchar(255) NOT NULL DEFAULT '@',
    `type`    enum('A','MX','CNAME','NS','SOA','PTR','TXT','AAAA','SVR','URL') NOT NULL,
    `data`    varchar(255) DEFAULT NULL,
    `ttl`     int(11) NOT NULL DEFAULT '3600',
    `mx`      int(11) DEFAULT NULL,
    `view`    varchar(32) NOT NULL DEFAULT 'any',
    `state`   varchar(20) NOT NULL DEFAULT 'running' COMMENT '状态',
    `remark`  varchar(100) NOT NULL DEFAULT '' COMMENT '描述',
    PRIMARY KEY (`id`),
    unique key (`zone_id`, `host`, `view`)
);


create table operation_log(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `operator` varchar(64) NOT NULL default  '' comment '修改的用户',
    `target_type` varchar(15) NOT NULL default  '' comment '记录类型',
    `target_id` bigint not null default 0 comment '记录id',
    `type` varchar(10) not null default '' comment '修改类型，add, update',
    `key_value` varchar(100) NOT NULL DEFAULT  '' comment '修改的记录',
    `diff` varchar(512) NOT NULL DEFAULT '' comment '修改值',
    primary key (`id`),
    key(target_type, target_id),
    key (`key_value`)
);