-- +goose Up
-- +goose StatementBegin

alter table user_banks add column
	client_id varchar(256) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table user_banks drop column
	client_id;
-- +goose StatementEnd

