-- +goose Up
-- +goose StatementBegin

create table users (
	user_email varchar(256) PRIMARY KEY NOT NULL
);
create table consents (
	consent_id varchar(256) PRIMARY KEY NOT NULL
);
create table banks (
	bank_id integer PRIMARY KEY NOT NULL,
	bank_name TEXT,
	bank_description TEXT,
	bank_api_url TEXT NOT NULL
);

create table user_banks (
	user_email varchar(256) NOT NULL REFERENCES users(user_email),
	bank_id INTEGER NOT NULL REFERENCES banks(bank_id),
	primary key (user_email, bank_id)
);

create table user_bank_consents (
	consent_id VARCHAR(256) NOT NULL REFERENCES consents(consent_id),
	bank_id INTEGER NOT NULL REFERENCES banks(bank_id),
	user_email varchar(256) NOT NULL REFERENCES users(user_email),
	primary key (consent_id,user_email, bank_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'nonsense';
-- +goose StatementEnd

