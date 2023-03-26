package main

import (
	"context"
	"files_storage/entity"
	"fmt"
	"log"
	"os"

	event_rabbitmq "files_storage/rabbitmq"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/file", func(c echo.Context) error {
		id := c.QueryParam("id")

		path := fmt.Sprintf("./questions/%s.txt", id)
		return c.File(path)
	})

	go func() {
		err := event_rabbitmq.ListenEvent(context.Background(), "new-question-file", func(req *entity.Question) {
			err := os.WriteFile(fmt.Sprintf("./questions/%s.txt", req.QuestionId), []byte(req.RunCode), 0644)
			if err != nil {
				log.Println("failed writing submitted code", err)
			}

			err = os.WriteFile(fmt.Sprintf("./questions/%s_answer.txt", req.QuestionId), []byte(req.StubCode), 0644)
			if err != nil {
				log.Println("failed writing submitted code", err)
			}
		})
		if err != nil {
			log.Println(err)
		}
	}()

	e.Logger.Fatal(e.Start(":8082"))

}
