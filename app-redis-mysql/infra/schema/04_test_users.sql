LOAD DATA INFILE '/names.txt' INTO TABLE names (name);

UPDATE names SET length = CHAR_LENGTH(name);