CREATE TABLE IF NOT EXISTS comments(
    comment_id SERIAL PRIMARY KEY,
    commenter INTEGER,
    req INTEGER,
    comment TEXT,
    created_ts DATE DEFAULT NOW(),
    CONSTRAINT comment_user_fk FOREIGN KEY (commenter)
        REFERENCES users (user_id) MATCH SIMPLE
        ON DELETE SET NULL
);

ALTER TABLE comments ADD CONSTRAINT comment_req_fk FOREIGN KEY (req)
REFERENCES requests (request_id) MATCH SIMPLE ON DELETE SET NULL;