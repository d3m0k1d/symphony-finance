# symphony-finance

# frontend repository

- https://github.com/nikitaLomeiko/vtb-api-hack_2025

# layout

we are doing our best to comply to [standard layout](https://github.com/golang-standards/project-layout)

- `./internal/api/uberproxy/` -- the main HTTP API package allowing authorized users to do requests to upstream APIs
- `./cmd/uberproxy/` -- executable uberproxy and composition root
- `./migrations/` -- sqlite database migrations for [goose](https://github.com/pressly/goose)
- `./client-pilot` -- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)-generated OpenBanking API Clients
	- build: `go generate ./client-pilot/`
- `./internal/client/` -- [hevaily WIP] API client optimized for ergonomics using generated clients.
	- `./internal/client/hack` -- implementation of interface with respect to the quirks of sandbox banks: missing and extra fields, nonstandard consents and auth
	- `./internal/client/gost` -- implementation strictly conforming to standards
- `./internal/multibank/` -- a multibank client for performing data
	aggregations over multiple banks. currently the only feature is combined
	transactions list
- `./internal/config/` -- configuration interface and providers backed by either environment variables or relational databases
- `./internal/otp/` -- service handling OTPs flow
- `./internal/mail/` -- email client for sending OTPs
	- `impl` -- real smtp client
	- `fake` -- mock logging codes, for development
# roadmap

## nov. 9

- proxying requests enriched with authentication
- combined transactions
- replaceable authentication and consent modules

## future

- support cloud-native config providers such as etcd
- gracefully handle combining different api versions
