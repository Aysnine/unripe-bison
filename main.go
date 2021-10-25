package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	requestId "github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"

	_ "github.com/Aysnine/unripe-bison/docs"
)

func main() {
	// Use an external setup function in order
	// to configure the app in tests as well
	app := Setup()

	// start the application
	// ! must with `localhost` on MacOS
	// ! https://medium.com/@leeprovoost/suppressing-accept-incoming-network-connections-warnings-on-osx-7665b33927ca
	log.Fatal(app.Listen("localhost:9000"))
}

// @title Server unripe-bison API
// @version 0.1.1
// @description Web server example by GoFiber⚡️ and CockroachDB📖
// @contact.name CNine
// @contact.email cnine229@gmail.com
// @license.name MIT
// @host unripe-bision.cnine.me
// @BasePath /
func Setup() *fiber.App {
	// connext database
	db := ConnectDB(getEnvVariable("DATABASE_CONNECTION"))

	// Initialize a new app
	app := fiber.New()

	// Web dashboard
	app.Get("/dashboard", monitor.New())

	// Swagger document
	app.Get("/swagger/*", swagger.Handler)

	// Default middleware config
	app.Use(requestId.New())

	// Default middleware config
	app.Use(logger.New())

	// Register the index route with a simple
	// "OK" response. It should return status
	// code 200
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Extract routing grouping
	SetupApi_Books(app, db)

	// Extract routing grouping
	SetupApi_HongKongWeather(app)

	// Return the configured app
	return app
}

// GetBooks godoc
// @Summary books
// @ID get-books
// @Produce  json
// @Router /api/books [get]
func SetupApi_Books(app *fiber.App, db *pgx.Conn) {
	// Routing grouping
	api := app.Group("/api")

	// Get books JSON response
	api.Get("/books", func(ctx *fiber.Ctx) error {

		// * Query

		rows, err := db.Query(context.Background(), "SELECT id, name FROM books")
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		defer rows.Close()

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

// GetHongKongWeather godoc
// @Summary hongkong weather info
// @ID get-hongkong-weather
// @Produce  json
// @Router /api/hongkong-weather [get]
func SetupApi_HongKongWeather(app *fiber.App) {
	// Routing grouping
	api := app.Group("/api")

	// Request other server
	api.Get("/hongkong-weather", func(ctx *fiber.Ctx) error {

		// * Request

		resp, err := http.Get("https://data.weather.gov.hk/weatherAPI/opendata/weather.php?dataType=fnd&lang=sc")
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}

		// * Marshaling

		type Response struct {
			GeneralSituation string `json:"generalSituation"`
		}

		response := Response{}
		json.Unmarshal(body, &response)

		// * Final response
		return ctx.JSON(response)
	})
}

func ConnectDB(connString string) *pgx.Conn {
	db, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = db.QueryRow(context.Background(), "select 'Database connected!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	return db
}

func getEnvVariable(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}
