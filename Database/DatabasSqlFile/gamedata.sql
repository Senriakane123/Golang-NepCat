/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80041
 Source Host           : localhost:3306
 Source Schema         : gamedata

 Target Server Type    : MySQL
 Target Server Version : 80041
 File Encoding         : 65001

 Date: 10/02/2025 18:07:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for allskilllist
-- ----------------------------
DROP TABLE IF EXISTS `allskilllist`;
CREATE TABLE `allskilllist`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `SkillName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `Des` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `PetID` int NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of allskilllist
-- ----------------------------

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
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of personalpetinfo
-- ----------------------------

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
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pet
-- ----------------------------
INSERT INTO `pet` VALUES (1, '炎龙', '火', '火焰喷射', 1.2, 1.5, 1, 1.3, 100, 50, 40, 80);
INSERT INTO `pet` VALUES (2, '冰狼', '冰', '寒冰爪击', 1.1, 1.4, 1.2, 1.1, 90, 55, 45, 85);
INSERT INTO `pet` VALUES (3, '雷鸟', '雷', '雷电冲击', 1.3, 1.2, 1.1, 1.4, 110, 45, 50, 90);

-- ----------------------------
-- Table structure for userinfo
-- ----------------------------
DROP TABLE IF EXISTS `userinfo`;
CREATE TABLE `userinfo`  (
  `ID` int NOT NULL AUTO_INCREMENT,
  `QQNum` int NOT NULL,
  `Name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `Item` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`ID`) USING BTREE,
  INDEX `QQNum`(`QQNum` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of userinfo
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
