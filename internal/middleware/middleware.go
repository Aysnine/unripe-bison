package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Will measure how long it takes before a response is returned
func ServerTiming() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()
		err := ctx.Next()
		stop := time.Now()

		// Do something with response
		ctx.Append("Server-Timing", fmt.Sprintf("app;duration=%v", stop.Sub(start).String()))

		// return stack error if exist
		return err
	}
}
