CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE(login)
);
