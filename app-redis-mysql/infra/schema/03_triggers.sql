USE words;

CREATE TRIGGER lower_case_trigger BEFORE INSERT ON indexedwords
FOR EACH ROW 
SET NEW.parsed_word = LOWER(NEW.word);

DELIMITER $$
CREATE TRIGGER extract_substrings_trigger
AFTER INSERT ON indexedwords
FOR EACH ROW
BEGIN
  DECLARE v_start INT DEFAULT 2;
  DECLARE v_string VARCHAR(255);

  SET v_string = NEW.parsed_word;
  WHILE v_start <= LENGTH(v_string) DO
    INSERT INTO subwords (subword, word)
    SELECT SUBSTRING(v_string, v_start), NEW.word;

    SET v_start = v_start + 1;
  END WHILE;
END $$
DELIMITER ;