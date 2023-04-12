package leaderboard_handler

import (
	"context"
	"io"
	"log"
	"server/entity"
	"server/entity/join_model"
	"server/handler"
	"server/helper"
	question_use_case "server/use_case/questions_use_case"
	"server/use_case/submission_use_case"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type LeaderboardHandlerImpl struct {
	Upgrader          websocket.Upgrader
	Rooms             map[string]*entity.GameRoom
	QuestionUseCase   question_use_case.QuestionUseCase
	SubmissionUseCase submission_use_case.SubmissionUseCase
	RedisConnection   *redis.Client
}

func NewLeaderboardHandler(qU question_use_case.QuestionUseCase, sU submission_use_case.SubmissionUseCase, rC *redis.Client) LeaderboardHandler {
	return &LeaderboardHandlerImpl{
		Upgrader:          websocket.Upgrader{},
		Rooms:             make(map[string]*entity.GameRoom),
		QuestionUseCase:   qU,
		SubmissionUseCase: sU,
		RedisConnection:   rC,
	}
}

func (h *LeaderboardHandlerImpl) CreateRoom(c echo.Context) error {
	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ws.Close()

	payload := new(entity.NewRoomPayload)
	for {
		err := ws.ReadJSON(payload)
		if err != nil {
			if err == io.EOF {
				continue
			}

			log.Println(err)
			break
		}

		payload.Player.PlayerId = "neoj"
		// payload.Player.PlayerId = helper.GenerateId(10)
		payload.Player.Connection = ws
		newRoom := &entity.GameRoom{
			// RoomId:   helper.GenerateId(20),
			RoomId:   "123",
			RoomName: payload.RoomName,
			MaxUser:  payload.MaxPlayer,
			Players:  []*entity.Player{payload.Player},
		}

		h.Rooms[newRoom.RoomId] = newRoom

		err = ws.WriteJSON(&entity.WSResponse{
			Message: handler.ROOM_CREATED,
			Data:    newRoom,
		})
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (h *LeaderboardHandlerImpl) JoinRoom(c echo.Context) error {
	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ws.Close()

	payload := new(entity.JoinRoomPayload)
	for {
		err := ws.ReadJSON(payload)
		if err != nil {
			if err == io.EOF {
				continue
			}

			log.Println(err)
			break
		}

		payload.Player.PlayerId = helper.GenerateId(10)
		payload.Player.Connection = ws
		v, found := h.Rooms[payload.RoomId]
		if !found {
			err := ws.WriteJSON(&entity.WSResponse{
				Message: handler.ROOM_NOT_FOUND,
			})
			if err != nil {
				log.Println(err)
			}
			break
		}

		if len(v.Players) == v.MaxUser {
			err := ws.WriteJSON(&entity.WSResponse{
				Message: handler.ROOM_FULL,
			})
			if err != nil {
				log.Println(err)
			}
			break
		}

		if v.IsGameStart {
			err := ws.WriteJSON(&entity.WSResponse{
				Message: handler.GAME_START,
			})
			if err != nil {
				log.Println(err)
			}
			break
		}

		v.Players = append(v.Players, payload.Player)
		h.Rooms[v.RoomId] = v
		for _, broadcast := range v.Players {
			if broadcast.Connection == ws {
				err := broadcast.Connection.WriteJSON(&entity.WSResponse{
					Message: handler.SUCCESS_JOIN_ROOM,
					Data:    v,
				})
				if err != nil {
					log.Println(err)
				}
			} else {
				err := broadcast.Connection.WriteJSON(&entity.WSResponse{
					Message: handler.SOMEONE_JOIN_ROOM,
					Data:    v,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}

		if len(v.Players) == v.MaxUser {
			for _, broadcast := range v.Players {
				err = broadcast.Connection.WriteJSON(&entity.WSResponse{
					Message: handler.ROOM_FULL_READY_TO_START,
					Data:    v,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil

}

func (h *LeaderboardHandlerImpl) StartGame(c echo.Context) error {
	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ws.Close()

	roomId := c.QueryParam("id")
	questions, err := h.QuestionUseCase.GetRandomQuestions(2)
	if err != nil {
		log.Println(err)
	}

	room, found := h.Rooms[roomId]
	if !found {
		err := ws.WriteJSON(&entity.WSResponse{
			Message: handler.ROOM_NOT_FOUND,
		})
		if err != nil {
			log.Println(err)
		}
	}

	room.Question = questions
	room.CurrentGameNumber++
	room.IsGameStart = true
	room.IsRoundStart = true

	h.Rooms[roomId] = room

	for _, player := range room.Players {
		err := player.Connection.WriteJSON(&entity.WSResponse{
			Message: handler.GAME_START,
			Data: entity.RoomAndQuestion{
				Room:     room,
				Question: room.Question[room.CurrentGameNumber-1].Question,
				StubCode: room.Question[room.CurrentGameNumber-1].StubCode,
			},
		})
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (h *LeaderboardHandlerImpl) SubmitCode2(c echo.Context) error {
	payload := new(entity.GameSubmissionPayload)
	err := c.Bind(payload)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(payload)

	return nil
}

func (h *LeaderboardHandlerImpl) SubmitCode(c echo.Context) error {
	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ws.Close()

	payload := new(entity.GameSubmissionPayload)
	err = ws.ReadJSON(payload)
	if err != nil {
		log.Println(err)
		return err
	}

	if payload.TestCaseLimit <= 0 {
		payload.TestCaseLimit = -1
	}

	log.Println("bind payload", payload)
	submissionId, err := h.SubmissionUseCase.NewSubmission(payload.SubmissionPayload, payload.TestCaseLimit)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("process", submissionId)
	submissionResult := new(join_model.SubmissionTestResult)
	for {
		res, err := h.SubmissionUseCase.GetSubmission(submissionId, payload.SubmissionPayload.CustomTestCase)
		if err != nil {
			log.Println(err)
			time.Sleep(70 * time.Millisecond)
			continue
		}

		if res.SubmissionStatus == "FINISHED" {
			submissionResult = res
			log.Println("process finished", res)
			break
		}

		time.Sleep(70 * time.Millisecond)
		log.Println("not yet")
	}

	playerConn := findPlayerWithId(payload.PlayerId, h.Rooms[payload.RoomId].Players)
	log.Println("check test case")
	for _, testResult := range submissionResult.TestResult {
		if testResult.Result != "PASS" {
			err := playerConn.Connection.WriteJSON(&entity.WSResponse{
				Message: handler.SOLUTION_REJECTED,
				Data:    submissionResult,
			})
			if err != nil {
				log.Println(err)
			}
		}
	}

	log.Println("write response")
	err = playerConn.Connection.WriteJSON(&entity.WSResponse{
		Message: handler.SOLUTION_ACCEPTED,
		Data:    submissionResult,
	})

	if err != nil {
		log.Println(err)
	}

	// tell everyone in the room
	room, found := h.Rooms[payload.RoomId]
	if !found {
		err := ws.WriteJSON(&entity.WSResponse{
			Message: handler.ROOM_NOT_FOUND,
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, player := range room.Players {
		if player.Connection != playerConn.Connection {
			err := ws.WriteJSON(&entity.WSResponse{
				Message: handler.SOMEONE_FINISH,
				Data:    player.DisplayName,
			})

			if err != nil {
				log.Println(err)
			}
		}
	}
	log.Println("done broadcast")

	// update leaderboard
	err = h.RedisConnection.ZIncrBy(context.Background(), payload.RoomId, 100, payload.PlayerId).Err()
	if err != nil {
		log.Println(err)
	}

	redisResult, err := h.RedisConnection.ZRevRangeWithScores(context.Background(), payload.RoomId, 0, -1).Result()
	if err != nil {
		log.Println(err)
	}

	for i, v := range redisResult {
		room.Leaderboard = append(room.Leaderboard, &entity.PlayerLeaderboard{
			PlayerId:    v.Member.(string),
			DisplayName: findPlayerWithId(v.Member.(string), room.Players).DisplayName,
			Score:       v.Score,
			Position:    i + 1,
		})
	}

	h.Rooms[payload.RoomId] = room
	log.Println("END PROCESS")
	return nil
}

func findPlayerWithId(playerId string, players []*entity.Player) *entity.Player {
	for _, player := range players {
		if player.PlayerId == playerId {
			return player
		}
	}

	return nil
}
