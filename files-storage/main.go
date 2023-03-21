package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/file", func(c echo.Context) error {
		id := c.QueryParam("id")

		path := fmt.Sprintf("./questions/%s/%s.js", id, id)
		fmt.Println(path)
		return c.File(path)
	})

	e.Logger.Fatal(e.Start(":8082"))
}
