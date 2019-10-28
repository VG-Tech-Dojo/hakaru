CREATE TABLE IF NOT EXISTS `eventlog` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `at` datetime DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `value` int(10) unsigned,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
