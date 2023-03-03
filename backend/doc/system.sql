create table `sys_user` (
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    `username`   varchar(64) NOT NULL DEFAULT '' COMMENT 'username',
    `password`   varchar(128) NOT NULL DEFAULT '' COMMENT 'password',
    `nickname`   varchar(128) NOT NULL DEFAULT '' COMMENT 'dns host',
    `status`     tinyint(2) NOT NULL DEFAULT '1' COMMENT 'status',
    PRIMARY KEY (`id`),
    unique key (username)
);

insert into sys_user(username, nickname, password, status) values ('admin', 'admin', '$2a$10$qTkz0WSwkYhv9RNLBttP7uCwnB2yGqgOLFIycsYn0S/GEATil2WwC', 1);