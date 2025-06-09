CREATE TABLE IF NOT EXISTS categories (
    id varchar(36) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    description char(255),
    is_active boolean NOT NULL
); 