CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGSERIAL NOT NULL,
    "order" VARCHAR(255) NOT NULL,
    sum NUMERIC NOT NULL,
    user_id BIGINT NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);
