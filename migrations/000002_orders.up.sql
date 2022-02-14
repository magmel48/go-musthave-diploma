CREATE TYPE OrderStatuses AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL NOT NULL,
    number VARCHAR(255) NOT NULL,
    status OrderStatuses NOT NULL DEFAULT 'NEW',
    user_id BIGINT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE(number),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);
