-- table 由程序自动创建 
-- create db gogo_admin
CREATE DATABASE gogo_admin DEFAULT character SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- create db gogo_main
CREATE DATABASE gogo_main DEFAULT character SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin`  (
      `id` bigint(20) NOT NULL AUTO_INCREMENT,
      `username` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
      `password` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
      `icon` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '头像',
      `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '邮箱',
      `nick_name` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '昵称',
      `note` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注信息',
      `login_time` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
      `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '帐号启用状态：0->禁用；1->启用',
      `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
      `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '后台用户表' ROW_FORMAT = DYNAMIC;


-- ----------------------------
-- Table structure for admin_role_relation
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_relation`;
CREATE TABLE `admin_role_relation`  (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `admin_id` bigint(20) NULL DEFAULT NULL,
    `role_id` bigint(20) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '后台用户和角色关系表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
     `id` bigint(20) NOT NULL AUTO_INCREMENT,
     `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '名称',
     `description` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '描述',
     `admin_count` int(11) NULL DEFAULT NULL COMMENT '后台用户数量',
     `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
     `status` int(1) NULL DEFAULT 1 COMMENT '启用状态：0->禁用；1->启用',
     `sort` int(11) NULL DEFAULT 0,
     PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '后台用户角色表' ROW_FORMAT = DYNAMIC;