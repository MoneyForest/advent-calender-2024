package usecase

import (
	"context"
	"os"

	"main/internal/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) CreateUser(ctx context.Context, email string) error {
	tracer := otel.Tracer(os.Getenv("SERVICE"))
	ctx, span := tracer.Start(ctx, "UserUsecase.CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("usecase", "UserUsecase"),
	)

	return u.userRepo.CreateUser(ctx, email)
}
