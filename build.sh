#!/bin/sh

go \
	build \
	-v \
	-o ./cmd/json2wasi2json-wazero-example/j2w2j \
	./cmd/json2wasi2json-wazero-example
