CREATE TABLE tasks (
    id char(36) PRIMARY KEY,
    family_id int,
    title char(50) NOT NULL,
    description char(255),
    is_recurring boolean DEFAULT FALSE,
    recurring_days json,
    created_by char(36),
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB;
