module github.com/go-api-libs/notion

go 1.27

require (
	cloud.google.com/go v0.123.0
	github.com/MarkRosemaker/cassette v0.0.0-20261011100210-a2f8d6b55f5a
	github.com/MarkRosemaker/errpath v0.0.0-20261011100139-0f973ba10436
	github.com/MarkRosemaker/jsonutil v0.0.0-20261011100139-a4e91ee48afc
	github.com/MarkRosemaker/openapi v0.0.0-20261011100231-67c4b5e30c66
	github.com/MarkRosemaker/openapi-codegen v0.0.0-20261011100643-ea735d1c5014
	github.com/MarkRosemaker/openapi-compress v0.0.0-20261011100434-3fc4653cc437
	github.com/MarkRosemaker/openapi-edit v0.0.0-20261011100309-209f26b6856c
	github.com/MarkRosemaker/openapi-enrich v0.0.0-20261011100357-62f6c0211d95
	github.com/MarkRosemaker/openapi-flatten v0.0.0-20261011100309-90ff3a922504
	github.com/ettle/strcase v0.2.0
	github.com/go-api-libs/api v0.0.0-20261011100139-4690447dcdfa
	github.com/go-api-libs/types v0.0.0-20260821232109-0cf45378823e
	golang.org/x/sync v0.24.0
)

require (
	github.com/MarkRosemaker/json2yaml v0.0.0-20261011100139-2d7c8ec80f68 // indirect
	github.com/MarkRosemaker/openapi-compare v0.0.0-20261011100308-25f7159f4d25 // indirect
	github.com/MarkRosemaker/openapi-merge v0.0.0-20261011100308-0f8ac90b88e7 // indirect
	github.com/MarkRosemaker/ordmap v0.0.0-20261011100159-2de6cfb102c4 // indirect
	github.com/MarkRosemaker/yaml v0.0.0-20261011100206-7a238ac9deda // indirect
	github.com/MarkRosemaker/yaml2json v0.0.0-20261011100139-1181ffd29042 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	golang.org/x/exp v0.0.0-20261011100332-faecfb24eb49 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/gofumpt v0.12.0 // indirect
)

tool (
	github.com/MarkRosemaker/openapi-codegen/cmd/openapi-codegen
	github.com/MarkRosemaker/openapi-compress/cmd/openapi-compress
	github.com/MarkRosemaker/openapi-enrich/cmd/openapi-enrich
	github.com/MarkRosemaker/openapi-flatten/cmd/openapi-flatten
)
