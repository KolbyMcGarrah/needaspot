ALTER TABLE requests DROP CONSTRAINT creator_fk;
ALTER TABLE requests ADD CONSTRAINT requests_creator_fk FOREIGN KEY (creator)
        REFERENCES users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION ON DELETE NO ACTION;