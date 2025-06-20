-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name TEXT UNIQUE NOT NULL,
                       description TEXT
);

CREATE TABLE permissions (
                             id SERIAL PRIMARY KEY,
                             code TEXT UNIQUE NOT NULL,
                             description TEXT
);

CREATE TABLE roles_permissions (
                                   role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
                                   permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
                                   PRIMARY KEY (role_id, permission_id)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       name TEXT,
                       role_id INTEGER REFERENCES roles(id) ON DELETE SET NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;