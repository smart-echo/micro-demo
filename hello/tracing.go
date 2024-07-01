package main

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

const (
	exporter = "OTEL_EXPORTER"
	protocol = "OTEL_EXPORTER_OTLP_PROTOCOL"
)

func initTracerProvider(service, version string) (*trace.TracerProvider, error) {
	ctx := context.Background()

	var exp trace.SpanExporter
	var err error

	expTyp := os.Getenv(exporter)
	switch expTyp {
	case "console", "":
		exp, err = stdouttrace.New()
	case "otlp":
		trace.NewTracerProvider()
		proto := os.Getenv(protocol)
		// ep := os.Getenv(endpoint)

		switch proto {
		case "http/protobuf", "":
			exp, err = otlptracehttp.New(ctx)
		case "grpc":
			exp, err = otlptracegrpc.New(ctx)
		default:
			err = fmt.Errorf("unsupported OTLP protocol: %s", protocol)
		}
	default:
		err = fmt.Errorf("unsupported exporter type: %s", expTyp)
	}

	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
