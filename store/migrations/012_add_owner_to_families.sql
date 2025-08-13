ALTER TABLE families 
ADD COLUMN owner_id CHAR(36),
ADD FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL;

---- create above / drop below ----

ALTER TABLE families
DROP FOREIGN KEY families_ibfk_1,
DROP COLUMN owner_id;