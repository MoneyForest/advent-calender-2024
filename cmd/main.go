package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	ddotel "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentelemetry"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/MoneyForest/timee-advent-calender-2024/internal"
	"github.com/MoneyForest/timee-advent-calender-2024/internal/handler"
)

type tracerProviderWrapper interface {
	Shutdown(context.Context) error
}

type ddTracerProvider struct {
	provider *ddotel.TracerProvider
}

func (d *ddTracerProvider) Shutdown(ctx context.Context) error {
	d.provider.Shutdown()
	return nil
}

func initTracer() (tracerProviderWrapper, error) {
	env := os.Getenv("ENV")
	if env == "dev" || env == "prod" {
		tp := ddotel.NewTracerProvider(
			ddtracer.WithService(os.Getenv("SERVICE")),
			ddtracer.WithEnv(os.Getenv("ENV")),
			ddtracer.WithServiceVersion("1.0.0"),
			ddtracer.WithAgentAddr("localhost:8126"),
		)
		otel.SetTracerProvider(tp)
		return &ddTracerProvider{provider: tp}, nil
	} else {
		exporter, err := otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithInsecure(),
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
		return tp, nil
	}
}

func main() {
	// Initialize tracer
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	// DI
	dependency := internal.DI()

	// Handlerの初期化
	userHandler := handler.NewUserHandler(dependency.UserUsecase)

	ctx := context.Background()
	b := make([]byte, 8)
	rand.Read(b)
	randomString := fmt.Sprintf("%x", b)

	if err := userHandler.CreateUser(ctx, "john.doe"+randomString+"@example.com"); err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
}
