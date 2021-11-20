package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Aysnine/unripe-bison/internal/database"
	"github.com/Aysnine/unripe-bison/internal/middleware"
	"github.com/Aysnine/unripe-bison/internal/service"
	"github.com/Aysnine/unripe-bison/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	requestId "github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
)

func init() {
	// Ready env
	godotenv.Load(".env")
}

func main() {
	// App ready timing
	start := time.Now()

	// connect chat redis
	chatRedis := database.ConnectRedis(os.Getenv("CHAT_REDIS_CONNECTION"))

	// Initialize a new app
	app := fiber.New()

	// Ready setup context
	setupContext := &types.SetupContext{App: app, ChatRedis: chatRedis}

	// App ready timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + "Env ready done!")

	// Use an external setup function in order
	// to configure the app in tests as well
	setupApp := Setup(setupContext)

	// start the application
	// ! must with `localhost` on MacOS
	// ! https://medium.com/@leeprovoost/suppressing-accept-incoming-network-connections-warnings-on-osx-7665b33927ca
	address := ":9000"
	if os.Getenv("MODE") == "local" {
		address = "localhost:9000"
	}
	log.Fatal(setupApp.Listen(address))
}

// @title UnripeBison Server API
// @version v0
// @description Web server example by GoFiber⚡️ and CockroachDB📖
// @contact.name Github
// @contact.url https://github.com/Aysnine/unripe-bison
// @license.name MIT
// @host unripe-bison.cnine.me
// @BasePath /
func Setup(setupContext *types.SetupContext) *fiber.App {
	// TODO move out when use mock
	// connect database
	db := database.ConnectPG(os.Getenv("DATABASE_CONNECTION"))
	setupContext.DB = db

	app := setupContext.App

	// App setup timing
	start := time.Now()

	if os.Getenv("MODE") == "local" || os.Getenv("MODE") == "development" {

		// Web monitor
		app.Get("/monitor", monitor.New())

		// Swagger document
		app.Static("/swagger/doc.json", "./docs/swagger.json")

		// Default middleware config
		app.Use(logger.New())

		// Custom Timing middleware
		app.Use(middleware.ServerTiming())
	}

	// Default middleware config
	app.Use(requestId.New())

	// Static Home Page
	app.Static("/", "./public")

	// Extract single route
	service.SetupApi_GetBooks(setupContext)
	service.SetupApi_AddBook(setupContext)
	service.SetupApi_GetHongKongWeather(setupContext)
	service.SetupApi_GetRandomAnimeImage(setupContext)
	service.SetupWebsocket_Chat(setupContext)

	// App setup timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + "All setup done!")

	// Return the configured app
	return app
}
