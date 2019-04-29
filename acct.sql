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
  `puser_id` varchar(36) DEFAULT '' COMMENT '第三方OPENID',
  `punion_id` varchar(36) DEFAULT '' COMMENT '第三方用户ID',
  `psession_key` varchar(50) DEFAULT '' COMMENT 'session_key',
  `access_token` varchar(128) DEFAULT '',
  `expires_in` int(11) DEFAULT '0',
  `jwt_token` varchar(256) DEFAULT '',
  `login_mode` tinyint(4) DEFAULT '0' COMMENT '登录方式：1：微信，2：手机注册',
  `login_name` varchar(50) DEFAULT '' COMMENT '登录账号',
  `login_pass` varchar(50) DEFAULT '' COMMENT '密码',
  `status` tinyint(4) DEFAULT '0' COMMENT '当前用户状态: 0： 正常 1：密码错误锁定,2：账户人工冻结（资金不可进出）,3 :账户止付（账户不允许消费）, 4：账户止 入（账户不允充值）',
  `avatar_url` varchar(256) DEFAULT '' COMMENT '用户头像URL',
  `nick_name` varchar(100) DEFAULT '' COMMENT '昵称',
  `mail` varchar(54) DEFAULT '' COMMENT '邮箱',
  `gender` tinyint(1) DEFAULT '0' COMMENT '性别',
  `phone` varchar(21) DEFAULT '' COMMENT '手机号码',
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
  KEY `idx_partner_user_id` (`puser_id`),
  KEY `idx_user_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8 COMMENT='账户信息表';