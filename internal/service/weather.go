package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
