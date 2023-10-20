package api

import "embed"

//go:embed openapi_ui/*
var UI embed.FS

//go:embed openapi.yaml
var OpenAPISpec embed.FS
