-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS userstable (
        user_id SERIAL PRIMARY KEY,
        login VARCHAR(20) NOT NULL,
        password VARCHAR(255) NOT NULL,
        token VARCHAR(255) NOT NULL,
        CONSTRAINT userstable_login_key UNIQUE (login)
    );

CREATE TABLE IF NOT EXISTS skinstable (
        skin_id SERIAL PRIMARY KEY,
		owner_name VARCHAR(20) NOT NULL,
        skin_name VARCHAR(30) NOT NULL,
        skin_type VARCHAR(10) NOT NULL,
        skin_src VARCHAR(255) NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS userstable;
DROP TABLE IF EXISTS skinstable;
-- +goose StatementEnd
