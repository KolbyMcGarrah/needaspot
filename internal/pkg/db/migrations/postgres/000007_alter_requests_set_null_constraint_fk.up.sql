ALTER TABLE requests DROP CONSTRAINT requests_creator_fk;
ALTER TABLE requests ADD CONSTRAINT creator_fk FOREIGN KEY (creator) REFERENCES
users(user_id) ON DELETE SET NULL;