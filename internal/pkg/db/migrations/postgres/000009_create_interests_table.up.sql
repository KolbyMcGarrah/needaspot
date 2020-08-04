CREATE TABLE IF NOT EXISTS interests(
    interest_id SERIAL PRIMARY KEY,
    user_id     INTEGER,
    request_id  INTEGER,
    accepted_user INTEGER,
    description TEXT,
    accepted BOOLEAN DEFAULT FALSE, 
    created_ts DATE DEFAULT NOW(),
    accepted_ts DATE,
    CONSTRAINT interest_user_id_fk FOREIGN KEY (user_id)
        REFERENCES users (user_id) MATCH SIMPLE
        ON DELETE SET NULL
);

ALTER TABLE interests ADD CONSTRAINT interest_request_id_fk FOREIGN KEY (request_id)
REFERENCES requests (request_id) MATCH SIMPLE ON DELETE SET NULL;
ALTER TABLE interests ADD CONSTRAINT interest_accepted_user_fk FOREIGN KEY (accepted_user)
REFERENCES users (user_id) MATCH SIMPLE ON DELETE SET NULL;