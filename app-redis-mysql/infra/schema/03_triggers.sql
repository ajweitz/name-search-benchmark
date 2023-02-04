USE words;

CREATE TRIGGER lower_case_trigger BEFORE INSERT ON indexedwords
FOR EACH ROW 
SET NEW.parsed_word = LOWER(NEW.word);

CREATE TRIGGER char_length_trigger BEFORE INSERT ON indexedwords
FOR EACH ROW 
SET NEW.length = CHAR_LENGTH(NEW.word);