CREATE TABLE families (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name char(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    is_active boolean DEFAULT TRUE
) ;
---- create above / drop below ----
drop table if exists families;