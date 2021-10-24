package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	// Use an external setup function in order
	// to configure the app in tests as well
	app := Setup()

	// Web dashboard
	app.Get("/dashboard", monitor.New())

	// start the application
	// ! must with `localhost` on MacOS
	// ! https://medium.com/@leeprovoost/suppressing-accept-incoming-network-connections-warnings-on-osx-7665b33927ca
	log.Fatal(app.Listen("localhost:9000"))
}

func Setup() *fiber.App {
	// connext database
	db := ConnectDB(getEnvVariable("DATABASE_CONNECTION"))

	// Initialize a new app
	app := fiber.New()

	// Register the index route with a simple
	// "OK" response. It should return status
	// code 200
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Routing grouping
	api := app.Group("/api")

	// Get all records from postgreSQL
	api.Get("/books", func(c *fiber.Ctx) error {
		// Select all book(s) from database
		rows, err := db.Query(context.Background(), "SELECT id, name FROM books order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		result := Books{}

		for rows.Next() {
			book := Book{}
			if err := rows.Scan(&book.ID, &book.Name); err != nil {
				return err // Exit if we get an error
			}

			// Append
			result.Books = append(result.Books, book)
		}

		// Return JSON
		return c.JSON(result)
	})

	// Return the configured app
	return app
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

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

type Book struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Books struct {
	Books []Book `json:"books"`
}
