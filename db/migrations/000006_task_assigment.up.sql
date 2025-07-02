CREATE TABLE task_assignments (
    id char(36) PRIMARY KEY,
    task_id char(36),
    assigned_to char(36),
    assigned_date date NOT NULL,
    status varchar(20) DEFAULT 'PENDING',
    accepted_by char(36),
    delegated_by char(36),
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_to) REFERENCES users(id),
    FOREIGN KEY (accepted_by) REFERENCES users(id),
    FOREIGN KEY (delegated_by) REFERENCES users(id)
) ENGINE=InnoDB;
