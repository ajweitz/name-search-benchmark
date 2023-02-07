USE words;


LOAD DATA INFILE '/var/lib/mysql-files/words.txt' INTO TABLE indexedwords (`word`);


INSERT INTO `words` SELECT * FROM `indexedwords`;