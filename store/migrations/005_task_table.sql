CREATE TABLE tasks (
    id CHAR(36) PRIMARY KEY,
    family_id INT,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_days JSON,
    status_id CHAR(36),
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES task_status(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);
---- create above / drop below ----
drop table IF EXISTS tasks;