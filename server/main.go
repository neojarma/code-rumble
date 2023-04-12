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
		return
	}

	rabbitConn, err := connection.GetConnectionRabbitMq()
	if err != nil {
		log.Println(err)
		return
	}

	redisConn, err := connection.GetConnectionRedis()
	if err != nil {
		log.Println(err)
		return
	}

	e := echo.New()
	router.Router(sqlConn, e, rabbitConn, redisConn)
	e.Logger.Fatal(e.Start(":8081"))
}
