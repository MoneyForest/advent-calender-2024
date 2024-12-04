package infrastructure

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	ddotel "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentelemetry"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ddTracerProviderWithDDOTel struct {
	provider *ddotel.TracerProvider
}

type ddTracerProvider struct {
	provider *sdktrace.TracerProvider
}

func InitDatadogWithDDOTel() (*ddTracerProviderWithDDOTel, error) {
	tp := ddotel.NewTracerProvider(
		ddtracer.WithService(os.Getenv("SERVICE")),
		ddtracer.WithEnv(os.Getenv("ENV")),
		ddtracer.WithServiceVersion("1.0.0"),
		ddtracer.WithAgentAddr("localhost:8126"),
	)
	otel.SetTracerProvider(tp)
	return &ddTracerProviderWithDDOTel{provider: tp}, nil
}

func (p *ddTracerProviderWithDDOTel) Shutdown(ctx context.Context) error {
	return p.provider.Shutdown()
}

func InitDatadog() (*ddTracerProvider, error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:4319"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithURLPath("/v1/traces"),
	)
	if err != nil {
		return nil, err
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
	return &ddTracerProvider{provider: tp}, nil
}

func (p *ddTracerProvider) Shutdown(ctx context.Context) error {
	return p.provider.Shutdown(ctx)
}
