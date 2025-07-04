CREATE TABLE sessions (
                          token CHAR(43) PRIMARY KEY,
                          data BLOB NOT NULL,
                          expiry TIMESTAMP(6) NOT NULL
) ENGINE =InnoDB;

CREATE INDEX sessions_expiry_idx ON sessions (expiry);