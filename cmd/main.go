package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/MoneyForest/timee-advent-calender-2024/internal"
	"github.com/MoneyForest/timee-advent-calender-2024/internal/handler"
	"github.com/MoneyForest/timee-advent-calender-2024/internal/infrastructure"
)

type TracerProviderWrapper interface {
	Shutdown(context.Context) error
}

func initTracer(env string) (TracerProviderWrapper, error) {
	switch env {
	case "dev":
		return infrastructure.InitDatadog()
	default:
		return infrastructure.InitJaeger()
	}
}

func main() {
	env := os.Getenv("ENV")
	// Initialize tracer
	tp, err := initTracer(env)
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
