USE words;

LOAD DATA INFILE '/var/lib/mysql-files/words.txt' INTO TABLE words (`word`);
LOAD DATA INFILE '/var/lib/mysql-files/lowcase.txt' INTO TABLE words (`parsed_word`);

LOAD DATA INFILE '/var/lib/mysql-files/words.txt' INTO TABLE indexedwords (`word`);
LOAD DATA INFILE '/var/lib/mysql-files/lowcase.txt' INTO TABLE indexedwords (`parsed_word`);

UPDATE words SET length = CHAR_LENGTH(`word`);
UPDATE indexedwords SET length = CHAR_LENGTH(`word`);
