package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Aysnine/unripe-bison/internal/database"
	"github.com/Aysnine/unripe-bison/internal/middleware"
	"github.com/Aysnine/unripe-bison/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	requestId "github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
)

func main() {
	// Use an external setup function in order
	// to configure the app in tests as well
	app := Setup()

	// start the application
	// ! must with `localhost` on MacOS
	// ! https://medium.com/@leeprovoost/suppressing-accept-incoming-network-connections-warnings-on-osx-7665b33927ca
	address := ":9000"
	if os.Getenv("MODE") == "development" {
		address = "localhost:9000"
	}
	log.Fatal(app.Listen(address))
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

	// Static Home Page
	app.Static("/", "./public")

	// Extract single route
	service.SetupApi_GetBooks(app, db)
	service.SetupApi_AddBook(app, db)
	service.SetupApi_GetHongKongWeather(app)
	service.SetupApi_GetRandomAnimeImage(app)

	// App setup timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + "All setup done!")

	// Return the configured app
	return app
}
