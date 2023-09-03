/*
 Navicat Premium Data Transfer

 Source Server         : test
 Source Server Type    : MySQL
 Source Server Version : 80025
 Source Host           : localhost:3306
 Source Schema         : douyin

 Target Server Type    : MySQL
 Target Server Version : 80025
 File Encoding         : 65001

 Date: 18/08/2023 20:20:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_info_id` bigint NULL DEFAULT NULL,
  `video_id` bigint NULL DEFAULT NULL,
  `content` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_user_infos_comments`(`user_info_id`) USING BTREE,
  INDEX `fk_videos_comments`(`video_id`) USING BTREE,
  CONSTRAINT `fk_user_infos_comments` FOREIGN KEY (`user_info_id`) REFERENCES `user_infos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_videos_comments` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comments
-- ----------------------------

-- ----------------------------
-- Table structure for user_favor_videos
-- ----------------------------
DROP TABLE IF EXISTS `user_favor_videos`;
CREATE TABLE `user_favor_videos`  (
  `user_info_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  PRIMARY KEY (`user_info_id`, `video_id`) USING BTREE,
  INDEX `fk_user_favor_videos_video`(`video_id`) USING BTREE,
  CONSTRAINT `fk_user_favor_videos_user_info` FOREIGN KEY (`user_info_id`) REFERENCES `user_infos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_user_favor_videos_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_favor_videos
-- ----------------------------

-- ----------------------------
-- Table structure for user_infos
-- ----------------------------
DROP TABLE IF EXISTS `user_infos`;
CREATE TABLE `user_infos`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `follow_count` bigint NULL DEFAULT NULL,
  `follower_count` bigint NULL DEFAULT NULL,
  `is_follow` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_infos
-- ----------------------------

-- ----------------------------
-- Table structure for user_logins
-- ----------------------------
DROP TABLE IF EXISTS `user_logins`;
CREATE TABLE `user_logins`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_info_id` bigint NULL DEFAULT NULL,
  `username` varchar(191) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `password` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`, `username`) USING BTREE,
  INDEX `fk_user_infos_user`(`user_info_id`) USING BTREE,
  CONSTRAINT `fk_user_infos_user` FOREIGN KEY (`user_info_id`) REFERENCES `user_infos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_logins
-- ----------------------------

-- ----------------------------
-- Table structure for user_relations
-- ----------------------------
DROP TABLE IF EXISTS `user_relations`;
CREATE TABLE `user_relations`  (
  `user_info_id` bigint NOT NULL,
  `follow_id` bigint NOT NULL,
  PRIMARY KEY (`user_info_id`, `follow_id`) USING BTREE,
  INDEX `fk_user_relations_follows`(`follow_id`) USING BTREE,
  CONSTRAINT `fk_user_relations_follows` FOREIGN KEY (`follow_id`) REFERENCES `user_infos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_user_relations_user_info` FOREIGN KEY (`user_info_id`) REFERENCES `user_infos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_relations
-- ----------------------------

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_info_id` bigint NULL DEFAULT NULL,
  `play_url` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `cover_url` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `favorite_count` bigint NULL DEFAULT NULL,
  `comment_count` bigint NULL DEFAULT NULL,
  `is_favorite` tinyint(1) NULL DEFAULT NULL,
  `title` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_user_infos_videos`(`user_info_id`) USING BTREE,
  CONSTRAINT `fk_user_infos_videos` FOREIGN KEY (`user_info_id`) REFERENCES `user_infos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of videos
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
