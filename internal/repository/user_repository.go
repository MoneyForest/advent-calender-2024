package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewUserRepository(db *sql.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		db:    db,
		redis: redis,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, email string) error {
	tracer := otel.Tracer(os.Getenv("SERVICE"))
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("repository", "UserRepository"),
		attribute.String("db", "mysql"),
	)

	// Begin transaction
	ctx, txSpan := tracer.Start(ctx, "mysql.begin_transaction")
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		txSpan.RecordError(err)
		txSpan.End()
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	txSpan.End()

	// Insert user
	ctx, insertSpan := tracer.Start(ctx, "mysql.execute_insert")
	insertSpan.SetAttributes(
		semconv.DBSystemMySQL,
		semconv.DBStatementKey.String("INSERT INTO users (email) VALUES (?)"),
	)

	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (email) VALUES (?)",
		email,
	)
	if err != nil {
		tx.Rollback()
		insertSpan.RecordError(err)
		insertSpan.End()
		return fmt.Errorf("failed to insert user: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	insertSpan.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
	insertSpan.End()

	// Commit transaction
	ctx, commitSpan := tracer.Start(ctx, "mysql.commit_transaction")
	if err := tx.Commit(); err != nil {
		commitSpan.RecordError(err)
		commitSpan.End()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	commitSpan.End()

	// Cache user in Redis
	ctx, redisSpan := tracer.Start(ctx, "redis.set_user")
	redisSpan.SetAttributes(
		attribute.String("redis.key", fmt.Sprintf("user:%s", email)),
		attribute.String("redis.operation", "SET"),
	)

	err = r.redis.Set(ctx, fmt.Sprintf("user:%s", email), time.Now().String(), 24*time.Hour).Err()
	if err != nil {
		redisSpan.RecordError(err)
		redisSpan.End()
		return fmt.Errorf("failed to cache user in Redis: %v", err)
	}
	redisSpan.End()

	return nil
}
