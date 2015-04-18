DROP TABLE IF EXISTS passport_login;
CREATE TABLE `passport_login` (
  `LoginID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `UserID` int(10) unsigned NOT NULL,
  `CreationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Source` smallint(5) unsigned NOT NULL DEFAULT '1',
  `LoginIp` bigint(20) unsigned NOT NULL,
  `AnonymousID` char(32) NOT NULL,
  `AuthCode` char(32) NOT NULL,
  `UserAgent` varchar(128) NOT NULL,
  PRIMARY KEY (`LoginID`),
  KEY `UserID` (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS passport_user;
CREATE TABLE `passport_user` (
  `UserID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `CreationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `BirthYear` year(4) NOT NULL,
  `Gender` enum('Secret','Male','Female') NOT NULL DEFAULT 'Secret',
  `Nickname` varchar(16) NOT NULL,
  PRIMARY KEY (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8;