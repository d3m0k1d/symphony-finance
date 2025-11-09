#!/usr/bin/env bash

oapi() {
	jq '
		(.. | objects | select(.anyOf) | select(.anyOf | length == 2)
		| select(.anyOf | any(.type == "null"))
		) |= (.type = (.anyOf[] | select(.type != "null") | .type) | .nullable = true | del(.anyOf))
	' ./client/openbanking-openapi.json > ./client/openbanking-openapi-fixed.json
}
oapi-paths() {
	in=$1
	if [[ $in =~ .*\.ya?ml ]]; then
		yaml=--yaml-input
	fi
		gojq $yaml -r '.paths|to_entries|.[].key' "$in" > "${in%.*}.txt"
}

bobgen() {
	# DON'T use any important db as such
	SQLITE_FILE=/tmp/hehehe.db
	rm -f $SQLITE_FILE
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(realpath $SQLITE_FILE) \
		go tool github.com/pressly/goose/v3/cmd/goose \
		-dir $(realpath ./migrations/) \
		up
	(
		export SQLITE_DSN=$(realpath $SQLITE_FILE)
		cd ./bobgen/
		go tool github.com/stephenafamo/bob/gen/bobgen-sqlite
	)
}

gen() {
	# oapi
	bobgen
	# go generate ./client-pilot/
}

build() {
	out=$1
	gen
	go build -o $1 ./cmd/uberproxy/
}


"$@"
