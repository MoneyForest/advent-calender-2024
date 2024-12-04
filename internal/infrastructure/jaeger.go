package infrastructure

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type jaegerTracerProvider struct {
	provider *sdktrace.TracerProvider
}

func InitJaeger() (jaegerTracerProvider, error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return jaegerTracerProvider{}, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("SERVICE")),
			semconv.DeploymentEnvironmentKey.String("local"),
		)),
	)
	otel.SetTracerProvider(tp)
	return jaegerTracerProvider{provider: tp}, nil
}

func (p jaegerTracerProvider) Shutdown(ctx context.Context) error {
	return p.provider.Shutdown(ctx)
}
