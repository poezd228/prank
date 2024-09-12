-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE tokens (
                        id BIGINT NOT NULL,
                        number BIGINT NOT NULL,
                        purpose INT NOT NULL,
                        secret TEXT NOT NULL,
                        expires_at BIGINT NOT NULL,
                        PRIMARY KEY (id, number, purpose)
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
