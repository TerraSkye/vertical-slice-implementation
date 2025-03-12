module github.com/terraskye/vertical-slice-implementation

go 1.22.0

toolchain go1.23.7

replace github.com/io-da/query => github.com/terraskye/query v0.0.0-20250310130952-cd3f17f32d38

require (
	github.com/go-kit/kit v0.13.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/io-da/query v1.3.5
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/trace v1.35.0
)

require (
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
)
