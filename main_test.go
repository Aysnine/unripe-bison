package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Aysnine/unripe-bison/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	// Define a structure for specifying input and output
	// data of a single test case. This structure is then used
	// to create a so called test map, which contains all test
	// cases, that should be run for testing this function
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "index route",
			route:         "/",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "",
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "Cannot GET /i-dont-exist",
		},
	}

	// Setup the app as it is done in the main function
	app := Setup(&types.SetupContext{
		App: fiber.New(),
	})

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		if test.expectedBody != "" {
			// Verify, that the response body equals the expected body
			assert.Equalf(t, test.expectedBody, string(body), test.description)
		}
	}
}

func TestApiBooksRoute(t *testing.T) {
	t.Parallel()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	app := Setup(&types.SetupContext{
		App: fiber.New(),
	})

	req, _ := http.NewRequest("GET", "/api/books", nil)
	res, err := app.Test(req, -1)

	assert.Equalf(t, false, err != nil, "error")
	assert.Equalf(t, 200, res.StatusCode, "error status code")
}

func TestHongKongWeatherRoute(t *testing.T) {
	app := Setup(&types.SetupContext{
		App: fiber.New(),
	})

	req, _ := http.NewRequest("GET", "/api/hongkong-weather", nil)
	res, err := app.Test(req, -1)

	assert.Equalf(t, false, err != nil, "error")
	assert.Equalf(t, 200, res.StatusCode, "error status code")
}

func TestRandomAnimeImageRoute(t *testing.T) {
	app := Setup(&types.SetupContext{
		App: fiber.New(),
	})

	req, _ := http.NewRequest("GET", "/api/random-anime-image", nil)
	res, err := app.Test(req, -1)

	assert.Equalf(t, false, err != nil, "error")
	assert.Equalf(t, 200, res.StatusCode, "error status code")
}
