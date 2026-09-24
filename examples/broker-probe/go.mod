module github.com/shing1211/webullapi4go/examples/broker-probe

go 1.26

require (
	github.com/shing1211/webullapi4go v0.0.0
	github.com/shing1211/webullapi4go/broker v0.0.0-20260921151224-4f0ce96cc403
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
)

replace github.com/shing1211/webullapi4go => ../../

replace github.com/shing1211/webullapi4go/broker => ../../broker
