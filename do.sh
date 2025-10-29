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

gen() {
	oapi
	go generate ./...
}

"$@"
