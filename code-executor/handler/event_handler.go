package handler

import "executor/entity"

type EventHandler interface {
	HandleEvent(event *entity.Submission)
}
