CREATE TABLE IF NOT EXISTS users
(
    id           INTEGER GENERATED ALWAYS AS IDENTITY,
    email        VARCHAR(100) UNIQUE,
    enc_password VARCHAR(512)
);

CREATE TABLE IF NOT EXISTS notes
(
    id         INTEGER GENERATED ALWAYS AS IDENTITY,
    email      VARCHAR(100)          NOT NULL REFERENCES users (email),
    title      VARCHAR(60)           NOT NULL,
    text       VARCHAR(2000),
    start_date TIMESTAMP             NOT NULL,
    end_date   TIMESTAMP,
    is_done    BOOLEAN DEFAULT FALSE NOT NULL
);

DO
$$
    BEGIN
        CREATE USER auth;
    EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;
    END
$$;
GRANT ALL PRIVILEGES ON users TO auth;

DO
$$
    BEGIN
        CREATE USER test;
    EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;
    END
$$;
GRANT ALL PRIVILEGES ON notes to test;
