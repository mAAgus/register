CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nick_name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    is_verificate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verification_token VARCHAR(64),
    token_created_at TIMESTAMP
);