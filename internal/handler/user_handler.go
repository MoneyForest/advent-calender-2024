package handler

import (
	"context"
	"os"

	"main/internal/usecase"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) CreateUser(ctx context.Context, email string) error {
	tracer := otel.Tracer(os.Getenv("SERVICE"))
	ctx, span := tracer.Start(ctx, "UserHandler.CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("handler", "UserHandler"),
		attribute.String("method", "CreateUser"),
	)

	return h.userUsecase.CreateUser(ctx, email)
}
