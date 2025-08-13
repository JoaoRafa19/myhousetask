CREATE TABLE
    sessions (
        token CHAR(43) PRIMARY KEY,
        data BLOB NOT NULL,
        expiry TIMESTAMP(6) NOT NULL
    ) ;

---- create above / drop below ----

DROP TABLE sessions;