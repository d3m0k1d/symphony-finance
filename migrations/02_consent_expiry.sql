-- +goose Up
-- +goose StatementBegin

alter table consents add column
	consent_expiry varchar(64) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table consents drop column
	consent_expiry;
-- +goose StatementEnd

