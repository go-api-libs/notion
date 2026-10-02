package notion

//go:generate go run api/fetch.go
//go:generate go tool openapi-enrich
//go:generate go tool openapi-flatten
//go:generate go tool openapi-compress
//go:generate go tool openapi-flatten
//go:generate go tool openapi-codegen -client -debug
