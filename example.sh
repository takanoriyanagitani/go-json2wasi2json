#!/bin/sh

jq -n -c '{helo:"wrld"}' |
	./cmd/json2wasi2json-wazero-example/j2w2j 
