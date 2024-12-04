package infrastructure

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	ddotel "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentelemetry"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ddTracerProvider struct {
	provider *ddotel.TracerProvider
}

func InitDatadog() (*ddTracerProvider, error) {
	tp := ddotel.NewTracerProvider(
		ddtracer.WithService(os.Getenv("SERVICE")),
		ddtracer.WithEnv(os.Getenv("ENV")),
		ddtracer.WithServiceVersion("1.0.0"),
		ddtracer.WithAgentAddr("localhost:8126"),
	)
	otel.SetTracerProvider(tp)
	return &ddTracerProvider{provider: tp}, nil
}

func (p *ddTracerProvider) Shutdown(ctx context.Context) error {
	return p.provider.Shutdown()
}
