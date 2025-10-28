package client

//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,client -package le -o le/gen.go accounts-v1.3.3le.yaml
//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,client -package pe -o pe/gen.go ./accounts-v1.3.7.yaml
////go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,client -package consent-pe -o consent-pe/gen.go account-consent-pe-2.1.yaml
