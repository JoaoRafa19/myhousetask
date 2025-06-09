CREATE TABLE family_invites (
    id char(36) PRIMARY KEY,
    family_id char(36),
    code char(255) UNIQUE NOT NULL,
    expires_at timestamp,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE
) ENGINE=InnoDB;
