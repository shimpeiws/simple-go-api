package main

import (
	"net/http"
	"io/ioutil"

	"github.com/labstack/echo"
)

func main() {
	// Echo instance
	e := echo.New()

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		// User ID from path `users/:id`
		id := c.Param("id")
		return c.String(http.StatusOK, id)
	})

	e.GET("/rails", func(c echo.Context) error {
		resp, err := http.Get("http://localhost:3000")
		if err != nil {
			return c.String(http.StatusNotFound, "Not Found")
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return c.String(http.StatusOK, string(body))
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
