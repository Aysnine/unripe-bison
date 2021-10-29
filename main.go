package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	requestId "github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
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

// @title UnripeBison Server API
// @version v0
// @description Web server example by GoFiber⚡️ and CockroachDB📖
// @contact.name Github
// @contact.url https://github.com/Aysnine/unripe-bison
// @license.name MIT
// @host unripe-bison.cnine.me
// @BasePath /
func Setup() *fiber.App {
	// App setup timing
	start := time.Now()

	// connext database
	db := ConnectDB(getEnvVariable("DATABASE_CONNECTION"))

	// Initialize a new app
	app := fiber.New()

	// Web monitor
	app.Get("/monitor", monitor.New())

	// Swagger document
	app.Static("/swagger/doc.json", "./docs/swagger.json")

	if getEnvVariable("MODE") == "development" {
		// Default middleware config
		app.Use(logger.New())
	}

	// Default middleware config
	app.Use(requestId.New())

	// Custom Timing middleware
	app.Use(ServerTiming())

	// Register the index route with a simple
	// "OK" response. It should return status
	// code 200
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Extract single route
	SetupApi_GetBooks(app, db)
	SetupApi_GetHongKongWeather(app)
	SetupApi_GetRandomAnimeImage(app)

	// App setup timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + "All setup done!")

	// Return the configured app
	return app
}

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

// GetHongKongWeather godoc
// @Summary hongkong weather info
// @ID get-hongkong-weather
// @Produce  json
// @Router /api/hongkong-weather [get]
func SetupApi_GetHongKongWeather(app *fiber.App) {
	// Request other server
	app.Get("/api/hongkong-weather", func(ctx *fiber.Ctx) error {

		// * Request

		start := time.Now()

		resp, err := http.Get("https://data.weather.gov.hk/weatherAPI/opendata/weather.php?dataType=fnd&lang=sc")
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		defer resp.Body.Close()

		stop := time.Now()
		ctx.Append("Server-Timing", fmt.Sprintf("rest;request=%v", stop.Sub(start).String()))

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

// GetRandomAnimeImage godoc
// @Summary hongkong weather info
// @ID get-random-anime-image
// @Produce  json
// @Router /api/random-anime-image [get]
func SetupApi_GetRandomAnimeImage(app *fiber.App) {
	// Request other server
	app.Get("/api/random-anime-image", func(ctx *fiber.Ctx) error {

		// * Request

		start := time.Now()

		resp, err := http.Get("https://api.waifu.pics/sfw/neko")
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		defer resp.Body.Close()

		stop := time.Now()
		ctx.Append("Server-Timing", fmt.Sprintf("rest;request=%v", stop.Sub(start).String()))

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}

		// * Unmarshal

		type Result struct {
			Url string `json:"url"`
		}

		result := Result{}
		json.Unmarshal(body, &result)

		// * Final response
		return ctx.Redirect(result.Url)
	})
}

func ConnectDB(connString string) *pgx.Conn {
	// Database connect timing
	start := time.Now()

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

	// Database connect timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + greeting)

	return db
}

func getEnvVariable(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}

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
