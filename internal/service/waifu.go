package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Aysnine/unripe-bison/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// GetRandomAnimeImage godoc
// @Summary hongkong weather info
// @ID get-random-anime-image
// @Produce  json
// @Router /api/random-anime-image [get]
func SetupApi_GetRandomAnimeImage(setup *types.SetupContext) {
	app := setup.App

	// Request other server
	app.Get("/api/random-anime-image/*", func(ctx *fiber.Ctx) error {

		// * Params
		rawPath := utils.ImmutableString(ctx.Params("*"))

		if len(rawPath) == 0 {
			rawPath = "sfw/neko"
		}

		// * Request

		restStart := time.Now()

		restResponse, restErr := http.Get("https://api.waifu.pics/" + rawPath)
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
