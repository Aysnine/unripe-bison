package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Aysnine/unripe-bison/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// GetBooks godoc
// @Summary books
// @ID get-books
// @Produce  json
// @Router /api/books [get]
func SetupApi_GetBooks(setup *types.SetupContext) {
	app := setup.App
	db := setup.DB

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
		ctx.Append("Server-Timing", fmt.Sprintf("sql;query=%v", stop.Sub(start).String()))

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

// GetBooks godoc
// @Summary books
// @ID add-book
// @Produce  json
// @Router /api/books [post]
func SetupApi_AddBook(setup *types.SetupContext) {
	app := setup.App
	db := setup.DB

	// Add a book
	app.Post("/api/books", func(ctx *fiber.Ctx) error {
		// * Params
		bookName := utils.ImmutableString(ctx.Query("name"))

		// * Insert

		start := time.Now()

		rows, err := db.Query(context.Background(), "INSERT INTO books(id, name) VALUES (gen_random_uuid(), $1) RETURNING id, name", bookName)
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		stop := time.Now()
		ctx.Append("Server-Timing", fmt.Sprintf("sql;insert=%v", stop.Sub(start).String()))

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
