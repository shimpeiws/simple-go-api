package main

import (
	"net/http"

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

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
