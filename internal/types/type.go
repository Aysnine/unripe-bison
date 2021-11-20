package types

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SetupContext struct {
	App       *fiber.App
	DB        *pgxpool.Pool
	ChatRedis *redis.Client
}
