CREATE TABLE IF NOT EXISTS requests(
    request_id SERIAL PRIMARY KEY,
    title VARCHAR (255) NOT NULL, 
    location VARCHAR (255) NOT NULL, 
    workout VARCHAR (255) NOT NULL,
    creator INTEGER NOT NULL,
    created_ts TIMESTAMP NOT NULL,
    CONSTRAINT requests_creator_fk FOREIGN KEY (creator)
        REFERENCES users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION ON DELETE NO ACTION
)