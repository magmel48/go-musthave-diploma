CREATE TABLE IF NOT EXISTS balances (
    id BIGSERIAL NOT NULL,
    current NUMERIC NOT NULL DEFAULT 0 CHECK (current >= 0),
    withdrawn NUMERIC NOT NULL DEFAULT 0 CHECK (withdrawn >= 0),
    user_id BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);
