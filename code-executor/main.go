package main

import (
	"context"
	"executor/connection"
	"executor/handler"
	event_rabbitmq "executor/rabbitmq"
	test_repository "executor/repository/test_result_repository"
	test_use_case "executor/use_case/test_result_use_case"
	"log"
	"server/repository/submission_repository"
)

func main() {
	sqlConn, err := connection.GetConnectionMySQL()
	if err != nil {
		log.Println(err)
		return
	}

	sr := submission_repository.NewSubmissionRepository(sqlConn)
	r := test_repository.NewTestResult(sqlConn)
	u := test_use_case.NewTestResult(r)
	h := handler.NewEventHandler(u, sr)

	ctx := context.Background()
	err = event_rabbitmq.ListenEvent(ctx, "submission", h.HandleEvent)
	if err != nil {
		log.Println(err)
		return
	}
}
