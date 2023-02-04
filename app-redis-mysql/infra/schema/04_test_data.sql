USE words;

LOAD DATA INFILE '/var/lib/mysql-files/words.txt' INTO TABLE words (`word`);

UPDATE words SET length = CHAR_LENGTH(`word`);