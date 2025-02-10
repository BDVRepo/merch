-- +goose Up
-- +goose StatementBegin
CREATE TABLE info_users (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE info_users;
-- +goose StatementEnd
