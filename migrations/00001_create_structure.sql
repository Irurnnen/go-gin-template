-- +goose Up
-- +goose StatementBegin
CREATE TYPE status AS ENUM (
    'available',
    'adopted',
    'sleeping',
    'playing',
    'sick'
);

CREATE TYPE gender AS ENUM (
    'male', 'female'
);

CREATE TABLE cats (
    id CHAR(36) NOT NULL UNIQUE PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    birth_timestamp TIMESTAMP,
    color VARCHAR(50),
    gender gender NOT NULL,
    status status NOT NULL,
    weight_kg INTEGER,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE toys (
    id CHAR(36) NOT NULL UNIQUE PRIMARY KEY,
    cat_id INTEGER NOT NULL REFERENCES cats(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50),
    color VARCHAR(50),
    material VARCHAR(50),
    created_at TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE toys;
DROP TABLE cats;
DROP ENUM gender;
DROP ENUM status;
-- +goose StatementEnd