-- +goose Up
CREATE TABLE waygate_links (
    id INTEGER PRIMARY KEY UNIQUE,
    waygate_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    link VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_waygate_links_waygate
    FOREIGN KEY (waygate_id)
    REFERENCES waygates(id)
    ON DELETE CASCADE
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
