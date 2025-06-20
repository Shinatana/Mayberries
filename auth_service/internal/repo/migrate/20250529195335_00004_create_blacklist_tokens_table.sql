-- +goose Up

CREATE TABLE token_blacklist (
                                 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                 token TEXT NOT NULL UNIQUE,
                                 blacklisted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

DROP TABLE IF EXISTS token_blacklist;
