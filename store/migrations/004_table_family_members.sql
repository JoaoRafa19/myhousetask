CREATE TABLE family_members (
    id char(36) PRIMARY KEY,
    family_id int,
    user_id char(36),
    role varchar(20) DEFAULT 'MEMBER',
    joined_at timestamp DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (family_id, user_id),
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ;

---- create above / drop below ----

drop table IF EXISTS family_members;