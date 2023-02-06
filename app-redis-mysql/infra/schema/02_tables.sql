USE words;

CREATE TABLE words (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(255) NOT NULL,
  `parsed_word` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE indexedwords (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(255) NOT NULL,
  `parsed_word` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX (`parsed_word`)
);

CREATE TABLE subwords (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `subword` varchar(255) NOT NULL,
  `word` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX (`subword`)
);