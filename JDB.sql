/*
Navicat MySQL Data Transfer

Source Server         : 10.193.1.19
Source Server Version : 50627
Source Host           : 10.193.1.19:3306
Source Database       : jcdb

Target Server Type    : MYSQL
Target Server Version : 50627
File Encoding         : 65001

Date: 2019-04-20 11:10:47
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `b_account`
-- ----------------------------
DROP TABLE IF EXISTS `b_account`;
CREATE TABLE `b_account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '登录账户ID',
  `wx_open_id` varchar(36) DEFAULT '' COMMENT '第三方OPENID',
  `wx_union_id` varchar(36) DEFAULT '' COMMENT '第三方用户ID',
  `wx_session_key` varchar(50) DEFAULT '' COMMENT 'session_key',
  `login_mode` tinyint(4) DEFAULT '0' COMMENT '登录方式：1：微信，2：手机注册',
  `login_name` varchar(50) DEFAULT '' COMMENT '登录账号',
  `login_pass` varchar(50) DEFAULT '' COMMENT '密码',
  `status` tinyint(4) DEFAULT '0' COMMENT '当前用户状态: 0： 正常 1：密码错误锁定,2：账户人工冻结（资金不可进出）,3 :账户止付（账户不允许消费）, 4：账户止 入（账户不允充值）',
  `avatar_url` varchar(100) DEFAULT NULL COMMENT '用户头像URL',
  `nick_name` varchar(100) DEFAULT '' COMMENT '昵称',
  `gender` tinyint(1) DEFAULT '0' COMMENT '性别',
  `city` varchar(100) DEFAULT '' COMMENT '城市',
  `province` varchar(100) DEFAULT '' COMMENT '省',
  `country` varchar(100) DEFAULT '' COMMENT '国家',
  `language` varchar(30) DEFAULT '' COMMENT '语言',
  `errors` tinyint(4) DEFAULT '0' COMMENT '密码错误次数',
  `account_bal` decimal(15,2) DEFAULT '0.00' COMMENT '账户余额-资金',
  `market` varchar(30) DEFAULT '' COMMENT '应用市场',
  `random_no` int(11) DEFAULT '0' COMMENT '用户随机数',
  `created_time` datetime DEFAULT NULL COMMENT '插入时间',
  `updated_time` datetime DEFAULT NULL COMMENT '更新时间',
  `memo` varchar(50) DEFAULT '' COMMENT '备注字段',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_partner_user_id` (`wx_open_id`),
  KEY `idx_user_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8 COMMENT='账户信息表';

-- ----------------------------
-- Records of b_account
-- ----------------------------

-- ----------------------------
-- Table structure for `b_disputes`
-- ----------------------------
DROP TABLE IF EXISTS `b_disputes`;
CREATE TABLE `b_disputes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '账户ID',
  `mct_no` varchar(10) DEFAULT '',
  `mct_trxn_no` varchar(50) DEFAULT '',
  `dispute_no` bigint(20) NOT NULL COMMENT '交易流水号',
  `dispute_date` datetime DEFAULT NULL COMMENT '交易时间',
  `dispute_amt` decimal(15,2) DEFAULT '0.00' COMMENT '交易金额',
  `status` char(1) DEFAULT '' COMMENT '交易处理状态 ',
  `dispute_memo` varchar(500) DEFAULT '' COMMENT '交易备注',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `update_user` varchar(50) DEFAULT '' COMMENT '人工调整交易的用户ID',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_trxn_date` (`dispute_date`),
  KEY `idx_proc_status` (`status`),
  KEY `idx_trxn_no` (`dispute_no`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='用户交易历史表';

-- ----------------------------
-- Records of b_disputes
-- ----------------------------

-- ----------------------------
-- Table structure for `b_flow`
-- ----------------------------
DROP TABLE IF EXISTS `b_flow`;
CREATE TABLE `b_flow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '账户ID',
  `mct_no` varchar(10) DEFAULT '',
  `mct_trxn_no` varchar(50) DEFAULT '',
  `trxn_no` bigint(20) NOT NULL COMMENT '交易流水号',
  `trxn_date` datetime DEFAULT NULL COMMENT '交易时间',
  `trxn_amt` decimal(15,2) DEFAULT '0.00' COMMENT '交易金额',
  `trxn_type` varchar(10) DEFAULT '' COMMENT '交易类型，包括资金交易，虚拟商品交易',
  `proc_status` char(1) DEFAULT '' COMMENT '交易处理状态 ',
  `proc_msg` varchar(10) DEFAULT '' COMMENT '交易处理结果原因',
  `account_bal` decimal(15,2) DEFAULT '0.00' COMMENT '账户余额',
  `trxn_memo` varchar(500) DEFAULT '' COMMENT '交易备注',
  `done_date` datetime DEFAULT NULL COMMENT '交易确认时间',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `update_user` varchar(50) DEFAULT '' COMMENT '人工调整交易的用户ID',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_trxn_date` (`trxn_date`),
  KEY `idx_proc_status` (`proc_status`),
  KEY `idx_trxn_no` (`trxn_no`)
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8 COMMENT='用户交易历史表';

-- ----------------------------
-- Records of b_flow
-- ----------------------------

-- ----------------------------
-- Table structure for `b_login`
-- ----------------------------
DROP TABLE IF EXISTS `b_login`;
CREATE TABLE `b_login` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '账户ID',
  `login_no` bigint(20) NOT NULL COMMENT '登录流水号',
  `login_time` datetime DEFAULT NULL COMMENT '登录时间',
  `login_result` tinyint(4) DEFAULT '0' COMMENT '登录结果',
  `duration` varchar(20) DEFAULT NULL,
  `device_ip` varchar(30) DEFAULT '' COMMENT '设备ip',
  `device_type` tinyint(4) DEFAULT '0' COMMENT '设备类型：1：ANDROID, 2：OS, 3：PC',
  `device_os` varchar(30) DEFAULT '' COMMENT '设备操作系统',
  `device_os_ver` varchar(30) DEFAULT '' COMMENT '设备操作系统版本',
  `device_id` varchar(30) DEFAULT '' COMMENT '设备id',
  `latitude` varchar(20) DEFAULT '' COMMENT '纬度',
  `longitude` varchar(20) DEFAULT '' COMMENT '经度',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_login_no` (`login_no`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COMMENT='用户登录账号表';

-- ----------------------------
-- Records of b_login
-- ----------------------------

DROP TABLE IF EXISTS `b_smscode`;
CREATE TABLE `b_smscode` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT NULL COMMENT '用户ID',
  `phone` varchar(21) DEFAULT '' COMMENT '手机号码',
  `sms_code` varchar(8) DEFAULT '' COMMENT '短信验证码',
  `sms_type` varchar(8) DEFAULT '' COMMENT '校验类型',
  `proc_status` varchar(5) DEFAULT '' COMMENT '发送结果代码',
  `proc_msg` char(50) DEFAULT '' COMMENT '发动结果原因',
  `status` char(1) DEFAULT 'e' COMMENT '默认代码',
  `verify_times` int(11) unsigned DEFAULT '0' COMMENT '校验次数',
  `valid_btime` datetime DEFAULT NULL COMMENT '有效开始时间',
  `valid_etime` datetime DEFAULT NULL COMMENT '有效结束时间',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `version` int(11) DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`),
  KEY `idx_user_phone` (`phone`) USING BTREE,
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=162 DEFAULT CHARSET=utf8 COMMENT='短信校验码记录表';