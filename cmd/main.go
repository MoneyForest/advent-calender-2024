package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"main/internal"
	"main/internal/handler"
	"main/internal/infrastructure"
)

func main() {
	// Initialize tracer
	tp, err := infrastructure.InitTracer(os.Getenv("ENV"))
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
	email := "john.doe+" + randomString + "@example.com"

	if err := userHandler.CreateUser(ctx, email); err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("success to create user: %v", email)
}
