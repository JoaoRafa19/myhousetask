CREATE TABLE calendar_events (
    id char(36) PRIMARY KEY,
    family_id int,
    task_id char(36),
    title char(50) NOT NULL,
    description char(255),
    start_time timestamp,
    end_time timestamp,
    created_by char(36),
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
) ;
---- create above / drop below ----
drop table if exists calendar_events;