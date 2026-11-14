module github.com/go-api-libs/notion

go 1.27

require (
	cloud.google.com/go v0.123.0
	github.com/MarkRosemaker/cassette v0.0.0-20261111155207-3bee9a17aeeb
	github.com/MarkRosemaker/errpath v0.0.0-20261109120451-91cd29644c73
	github.com/MarkRosemaker/jsonutil v0.0.0-20261114002746-7598c9df8b8f
	github.com/MarkRosemaker/openapi v0.0.0-20261114002923-49c60d62ff12
	github.com/MarkRosemaker/openapi-codegen v0.0.0-20261112182256-79f0fa399420
	github.com/MarkRosemaker/openapi-compress v0.0.0-20261114035644-54d90cd45ea4
	github.com/MarkRosemaker/openapi-edit v0.0.0-20261114003129-0b34f69d9d92
	github.com/MarkRosemaker/openapi-enrich v0.0.0-20261112181646-990112080a63
	github.com/MarkRosemaker/openapi-flatten v0.0.0-20261112181512-46b453b2a7e8
	github.com/ettle/strcase v0.2.0
	github.com/go-api-libs/api v0.0.0-20261112184314-5e0f94e7887c
	github.com/go-api-libs/types v0.0.0-20261109120451-6cd1daca4971
	golang.org/x/sync v0.25.0
)

require (
	github.com/MarkRosemaker/json2yaml v0.0.0-20261114002559-5de2992f679e // indirect
	github.com/MarkRosemaker/openapi-compare v0.0.0-20261114003127-a7b12b03c505 // indirect
	github.com/MarkRosemaker/openapi-merge v0.0.0-20261112181510-137dd05508c2 // indirect
	github.com/MarkRosemaker/ordmap v0.0.0-20261113130919-027ff7850b50 // indirect
	github.com/MarkRosemaker/yaml v0.0.0-20261114002825-bdb424d364db // indirect
	github.com/MarkRosemaker/yaml2json v0.0.0-20261114002559-cb6e42b758d5 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	golang.org/x/exp v0.0.0-20261113113407-e3970619a1a4 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	golang.org/x/tools v0.52.1-0.20261106145633-d407941a339a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/gofumpt v0.12.0 // indirect
)

tool (
	github.com/MarkRosemaker/openapi-codegen/cmd/openapi-codegen
	github.com/MarkRosemaker/openapi-compress/cmd/openapi-compress
	github.com/MarkRosemaker/openapi-enrich/cmd/openapi-enrich
	github.com/MarkRosemaker/openapi-flatten/cmd/openapi-flatten
)
