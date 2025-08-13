
CREATE INDEX sessions_expiry_idx ON sessions (expiry);

---- create above / drop below ----

DROP INDEX sessions_expiry_idx ON sessions;