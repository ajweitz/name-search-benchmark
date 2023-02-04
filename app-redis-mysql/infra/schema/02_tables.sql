USE words;

CREATE TABLE words (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(255) NOT NULL,
  `length` int(10) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX (`length`,`word`)
);