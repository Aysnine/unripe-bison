package types

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SetupContext struct {
	App       *fiber.App
	DB        *pgxpool.Pool
	ChatStore *redis.Storage
}
