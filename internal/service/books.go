package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

// GetBooks godoc
// @Summary books
// @ID get-books
// @Produce  json
// @Router /api/books [get]
func SetupApi_GetBooks(app *fiber.App, db *pgx.Conn) {
	// Get books JSON response
	app.Get("/api/books", func(ctx *fiber.Ctx) error {

		// * Query

		start := time.Now()

		rows, err := db.Query(context.Background(), "SELECT id, name FROM books")
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		stop := time.Now()
		ctx.Append("Server-Timing", fmt.Sprintf("db;query=%v", stop.Sub(start).String()))

		// * Marshaling

		type Book struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		type Response struct {
			Books []Book `json:"books"`
		}

		response := Response{Books: make([]Book, 0)}

		for rows.Next() {
			book := Book{}
			if err := rows.Scan(&book.ID, &book.Name); err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			response.Books = append(response.Books, book)
		}

		// * Final response
		return ctx.JSON(response)
	})
}
