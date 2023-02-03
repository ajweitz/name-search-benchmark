USE words;

CREATE TABLE names (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `length` int(10) NOT NULL
  PRIMARY KEY (`id`)
  INDEX (`length`)
);