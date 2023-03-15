package main

import (
	"log"
	"server/connection"
	"server/router"

	"github.com/labstack/echo/v4"
)

func main() {
	sqlConn, err := connection.GetConnectionMySQL()
	if err != nil {
		log.Println(err)
	}

	e := echo.New()
	router.Router(sqlConn, e)
	e.Logger.Fatal(e.Start(":8080"))
}