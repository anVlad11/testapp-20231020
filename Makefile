#!/usr/bin/make

-include include.make

run-local:
	go run ./cmd/app/ --config-path ./config.yaml

generate-oapi-models:
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@d516da7
	@mkdir -p ./pkg/testapp
	@oapi-codegen --config oapi-codegen.yaml docs/api/openapi.yaml

generate: generate-oapi-models