CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL, 
    password VARCHAR(127) NOT NULL, 
    age INT NOT NULL, 
    gender VARCHAR(10) NOT NULL, 
    level VARCHAR(10) NOT NULL,
    created_on TIMESTAMP NOT NULL, 
    last_login TIMESTAMP
)