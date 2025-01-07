Drop database IF EXISTS OPChainA;

Create Database If Not Exists OPChainA Character Set UTF8;
USE OPChainA;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET GLOBAL binlog_format = 'ROW';


Drop database IF EXISTS OPChainB;

Create Database If Not Exists OPChainB Character Set UTF8;
USE OPChainB;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET GLOBAL binlog_format = 'ROW';