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

	"github.com/Aysnine/unripe-bison/internal/database"
	"github.com/Aysnine/unripe-bison/internal/middleware"

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
	// Ready env
	godotenv.Load(".env")

	// App setup timing
	start := time.Now()

	// connext database
	db := database.ConnectPG(os.Getenv("DATABASE_CONNECTION"))

	// Initialize a new app
	app := fiber.New()

	// Web monitor
	app.Get("/monitor", monitor.New())

	// Swagger document
	app.Static("/swagger/doc.json", "./docs/swagger.json")

	if os.Getenv("MODE") == "development" {
		// Default middleware config
		app.Use(logger.New())
	}

	// Default middleware config
	app.Use(requestId.New())

	// Custom Timing middleware
	app.Use(middleware.ServerTiming())

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

		restStart := time.Now()

		restResponse, restErr := http.Get("https://api.waifu.pics/sfw/neko")
		if restErr != nil {
			return ctx.Status(500).SendString(restErr.Error())
		}
		defer restResponse.Body.Close()

		restStop := time.Now()
		ctx.Append("Server-Timing", fmt.Sprintf("rest;request=%v", restStop.Sub(restStart).String()))

		restBody, restErr := io.ReadAll(restResponse.Body)
		if restErr != nil {
			return ctx.Status(500).SendString(restErr.Error())
		}

		// * Rest Marshaling

		type RestResult struct {
			Url string `json:"url"`
		}

		restResult := RestResult{}
		json.Unmarshal(restBody, &restResult)

		// * File stream read and write

		fileRestResponse, fileRestErr := http.Get(restResult.Url)
		if fileRestErr != nil {
			return ctx.Status(500).SendString(fileRestErr.Error())
		}
		defer fileRestResponse.Body.Close()

		ctx.SendStream(fileRestResponse.Body)

		return nil
	})
}
