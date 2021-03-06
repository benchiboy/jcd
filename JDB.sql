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
CREATE TABLE `b_flow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '账户ID',
  `mct_no` varchar(10) DEFAULT '',
  `mct_trxn_no` varchar(50) DEFAULT '',
  `trxn_no` bigint(20) NOT NULL COMMENT '交易流水号',
  `trxn_amt` int(11) unsigned DEFAULT '0' COMMENT '交易金额',
  `trxn_type` varchar(10) DEFAULT '' COMMENT '交易类型，包括资金交易，虚拟商品交易',
  `proc_status` char(1) DEFAULT '' COMMENT '交易处理状态 ',
  `proc_msg` varchar(10) DEFAULT '' COMMENT '交易处理结果原因',
  `account_bal` decimal(15,2) DEFAULT '0.00' COMMENT '账户余额',
  `code_url` varchar(32) DEFAULT '',
  `prpay_id` varchar(32) DEFAULT '',
  `trxn_memo` varchar(500) DEFAULT '' COMMENT '交易备注',
  `trxn_date` datetime DEFAULT NULL COMMENT '交易时间',
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
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8 COMMENT='用户交易历史表';
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


DROP TABLE IF EXISTS `b_index`;
CREATE TABLE `b_index` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `index_subject` varchar(50) DEFAULT NULL COMMENT '指标主题',
  `index_category` varchar(50) DEFAULT NULL COMMENT '指标种类',
  `index_name` varchar(50) DEFAULT '0' COMMENT '指标名称',
  `index_value` decimal(10,2) DEFAULT NULL COMMENT '指标值',
  `index_date` int(11) DEFAULT NULL COMMENT '指标日期',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='观察指标表';


DROP TABLE IF EXISTS `b_badloan`;
CREATE TABLE `b_badloan` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mct_no` varchar(10) DEFAULT '' COMMENT '商户编号',
  `loan_no` varchar(50) DEFAULT '' COMMENT '合同号',
  `user_name` varchar(50) DEFAULT '' COMMENT '用户姓名',
  `id_no` varchar(21) DEFAULT '' COMMENT '证件号',
  `addr` varchar(250) DEFAULT '' COMMENT '住址',
  `gender` varchar(2) DEFAULT '' COMMENT '性别',
  `phone` varchar(21) DEFAULT '' COMMENT '电话',
  `memo` varchar(1024) NOT NULL DEFAULT '' COMMENT '备注',
  `status` char(1) DEFAULT '' COMMENT '状态',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `update_user` varchar(50) DEFAULT '' COMMENT '人工调整交易的用户ID',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='坏账用户表';


DROP TABLE IF EXISTS `b_comment_user`;
CREATE TABLE `b_comment_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `action_type` tinyint(4) DEFAULT NULL COMMENT '动作类型',
  `comm_no` bigint(20) DEFAULT '0' COMMENT '评论编号',
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '账户ID',
  `insert_time` datetime DEFAULT NULL COMMENT '插入时间',
  `version` int(11) unsigned DEFAULT '0' COMMENT '版本',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='评论对应的用户表';