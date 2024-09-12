-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE public.users
(
    id             BIGSERIAL,
    login          VARCHAR(255) PRIMARY KEY,
    password       VARCHAR(255) NOT NULL,
    name           VARCHAR(255),
    role           VARCHAR(255),
    login_count    INT,
    premium      BOOLEAN,
    email VARCHAR(255),
    confirmed BOOLEAN,
    telegram_id BIGINT
);



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE IF NOT EXIST public.users
-- +goose StatementEnd
