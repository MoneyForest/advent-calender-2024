package internal

import (
	"context"
	"database/sql"
	"log"

	"github.com/MoneyForest/timee-advent-calender-2024/internal/repository"
	"github.com/MoneyForest/timee-advent-calender-2024/internal/usecase"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

type Dependency struct {
	UserUsecase *usecase.UserUsecase
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/maindb?parseTime=true")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func initRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func DI() *Dependency {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	cache, err := initRedis()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db, cache)
	userUsecase := usecase.NewUserUsecase(userRepo)

	return &Dependency{
		UserUsecase: userUsecase,
	}
}
