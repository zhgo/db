CREATE TABLE `table1` (
  `UserID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '',
  `CreationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
  `BirthYear` year(4) NOT NULL COMMENT '',
  `Gender` enum('Secret','Male','Female') NOT NULL DEFAULT 'Secret' COMMENT '',
  `Nickname` varchar(16) NOT NULL COMMENT '',
  PRIMARY KEY (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8 COMMENT='';