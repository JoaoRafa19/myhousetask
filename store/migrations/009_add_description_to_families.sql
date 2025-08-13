
ALTER TABLE families ADD COLUMN description TEXT;

CREATE TABLE task_status (
                             id CHAR(36) PRIMARY KEY,
                             name VARCHAR(20) NOT NULL UNIQUE,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----

ALTER TABLE families DROP COLUMN description;
DROP TABLE task_status;