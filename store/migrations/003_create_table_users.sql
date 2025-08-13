CREATE TABLE users (
    id char(36) PRIMARY KEY,
    name char(255) NOT NULL,
    email char(255) UNIQUE NOT NULL,
    password_hash char(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
) ;
---- create above / drop below ----
drop table if exists users;