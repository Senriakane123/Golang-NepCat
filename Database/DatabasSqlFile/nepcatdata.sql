/*
 Navicat Premium Data Transfer

 Source Server         : msqDb
 Source Server Type    : MySQL
 Source Server Version : 80041
 Source Host           : localhost:3306
 Source Schema         : nepcatdata

 Target Server Type    : MySQL
 Target Server Version : 80041
 File Encoding         : 65001

 Date: 22/02/2025 09:20:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for adminuser
-- ----------------------------
DROP TABLE IF EXISTS `adminuser`;
CREATE TABLE `adminuser`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `QQNum` int NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of adminuser
-- ----------------------------
INSERT INTO `adminuser` VALUES (1, 735439479);

-- ----------------------------
-- Table structure for allskilllist
-- ----------------------------
DROP TABLE IF EXISTS `allskilllist`;
CREATE TABLE `allskilllist`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `SkillName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `Des` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `PetID` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of allskilllist
-- ----------------------------
INSERT INTO `allskilllist` VALUES (1, '普通攻击', '最简单的攻击方式，没有属性变化', '“”');

-- ----------------------------
-- Table structure for personalpetinfo
-- ----------------------------
DROP TABLE IF EXISTS `personalpetinfo`;
CREATE TABLE `personalpetinfo`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `UserID` int NOT NULL,
  `PetId` int NOT NULL,
  `PetLevel` int NOT NULL,
  `Exp` int NOT NULL,
  `QQNum` int NOT NULL,
  `Skill` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE,
  INDEX `UserID`(`UserID` ASC) USING BTREE,
  INDEX `PetId`(`PetId` ASC) USING BTREE,
  INDEX `personalpetinfo_ibfk_3`(`QQNum` ASC) USING BTREE,
  CONSTRAINT `personalpetinfo_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `userinfo` (`ID`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `personalpetinfo_ibfk_2` FOREIGN KEY (`PetId`) REFERENCES `pet` (`ID`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `personalpetinfo_ibfk_3` FOREIGN KEY (`QQNum`) REFERENCES `userinfo` (`QQNum`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of personalpetinfo
-- ----------------------------
INSERT INTO `personalpetinfo` VALUES (11, 11, 1, 1, 0, 735439479, '[{\"ID\":1,\"SkillName\":\"\",\"Des\":\"\",\"PetID\":\"\"}]');

-- ----------------------------
-- Table structure for pet
-- ----------------------------
DROP TABLE IF EXISTS `pet`;
CREATE TABLE `pet`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `Type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `Skill` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `HealthGrowthFactor` float NULL DEFAULT NULL,
  `AtkGrowthFactor` float NULL DEFAULT NULL,
  `DefenseGrowthFactor` float NULL DEFAULT NULL,
  `EnergyGrowthFactor` float NULL DEFAULT NULL,
  `BaseHealth` int NOT NULL,
  `BaseAtk` int NOT NULL,
  `BaseDef` int NOT NULL,
  `BaseEnergy` int NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of pet
-- ----------------------------
INSERT INTO `pet` VALUES (1, '炎龙', '火', '火焰喷射', 1.2, 1.5, 1, 1.3, 100, 50, 40, 80);
INSERT INTO `pet` VALUES (2, '冰狼', '冰', '寒冰爪击', 1.1, 1.4, 1.2, 1.1, 90, 55, 45, 85);
INSERT INTO `pet` VALUES (3, '雷鸟', '雷', '雷电冲击', 1.3, 1.2, 1.1, 1.4, 110, 45, 50, 90);

-- ----------------------------
-- Table structure for rgsgroup
-- ----------------------------
DROP TABLE IF EXISTS `rgsgroup`;
CREATE TABLE `rgsgroup`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `GroupID` int NOT NULL,
  `SeessionID` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of rgsgroup
-- ----------------------------
INSERT INTO `rgsgroup` VALUES (5, 1016759932, 'e6132222d6b445ed9b026b7272e29c96');
INSERT INTO `rgsgroup` VALUES (6, 907318559, 'ae074a988cd444adad5bac27f3b3f991');
INSERT INTO `rgsgroup` VALUES (7, 915116547, '96932de066d447648a8b0abf8df20bf7');

-- ----------------------------
-- Table structure for userinfo
-- ----------------------------
DROP TABLE IF EXISTS `userinfo`;
CREATE TABLE `userinfo`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `QQNum` int NOT NULL,
  `Name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `Item` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `SignInDayCout` int NOT NULL,
  `SignInTime` datetime NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE,
  INDEX `QQNum`(`QQNum` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of userinfo
-- ----------------------------
INSERT INTO `userinfo` VALUES (11, 735439479, 'Falling under the seabed', '{\"1\":2,\"2\":10}', 1, '2025-02-20 15:55:33');

SET FOREIGN_KEY_CHECKS = 1;
