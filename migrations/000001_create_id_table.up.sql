CREATE TABLE IF NOT EXISTS `T_ID` (
  `id` VARCHAR (36) NOT NULL PRIMARY KEY,
  `total` TINYINT UNSIGNED NOT NULL DEFAULT '0',
  `is_end` BIT(1) NOT NULL DEFAULT b'0',
  `expire` VARCHAR(19) NOT NULL DEFAULT '1970-01-01T09:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
