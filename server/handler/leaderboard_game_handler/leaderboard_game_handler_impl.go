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
		err := h.readJSON(ws, payload)
		if err != nil {
			h.writeErrorResponse(ws, handler.FAILED_TO_CREATE_ROOM, handler.INVALID_BODY_REQUEST)
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
		h.writeSuccessResponse(ws, handler.ROOM_CREATED, newRoom)
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
		err := h.readJSON(ws, payload)
		if err != nil {
			h.writeErrorResponse(ws, handler.FAILED_TO_CREATE_ROOM, handler.INVALID_BODY_REQUEST)
			break
		}

		payload.Player.PlayerId = helper.GenerateId(10)
		payload.Player.Connection = ws
		gameRoom, err := h.isRoomExist(ws, payload.RoomId, h.Rooms)
		if err != nil {
			h.writeErrorResponse(ws, handler.FAILED_TO_JOIN_ROOM, handler.ROOM_IS_NOT_EXIST)
			break
		}

		if h.isRoomFull(ws, gameRoom) {
			h.writeErrorResponse(ws, handler.FAILED_TO_JOIN_ROOM, handler.ROOM_IS_FULL)
			break
		}

		if h.isRoomAlreadyStart(ws, gameRoom) {
			h.writeErrorResponse(ws, handler.FAILED_TO_JOIN_ROOM, handler.ROOM_IS_ALREADY_PLAYED)
			break
		}

		h.Rooms[gameRoom.RoomId] = h.addNewPlayer(gameRoom, payload.Player)
		h.broadcastMessageWhenSomeoneJoinRoom(ws, gameRoom)
		h.broadcastMessageWhenRoomReadyToStart(gameRoom)
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
	gameRoom, err := h.isRoomExist(ws, roomId, h.Rooms)
	if err != nil {
		h.writeErrorResponse(ws, handler.FAILED_TO_JOIN_ROOM, handler.ROOM_IS_NOT_EXIST)
		return err
	}

	questions, err := h.QuestionUseCase.GetRandomQuestions(2)
	if err != nil {
		log.Println(err)
	}

	gameRoom.Question = questions
	gameRoom.CurrentGameNumber++
	gameRoom.IsGameStart = true
	gameRoom.IsRoundStart = true

	h.Rooms[roomId] = gameRoom

	h.broadcastMessageWhenGameIsStart(gameRoom)

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
	err = h.readJSON(ws, payload)
	if err != nil {
		h.writeErrorResponse(ws, handler.INVALID_BODY_REQUEST, err)
		return err
	}

	if payload.TestCaseLimit <= 0 {
		payload.TestCaseLimit = -1
	}

	submissionId, err := h.SubmissionUseCase.NewSubmission(payload.SubmissionPayload, payload.TestCaseLimit)
	if err != nil {
		log.Println(err)
		return err
	}

	submissionResult := h.getSubmissionResult(submissionId, payload.SubmissionPayload.CustomTestCase)

	playerReq := h.findPlayerWithId(payload.PlayerId, h.Rooms[payload.RoomId].Players)
	if h.isAllTestCasesPassed(playerReq.Connection, submissionResult) {
		err = playerReq.Connection.WriteJSON(&entity.WSResponse{
			Message: handler.SOLUTION_ACCEPTED,
			Data:    submissionResult,
		})

		if err != nil {
			log.Println(err)
		}
	} else {
		err := playerReq.Connection.WriteJSON(&entity.WSResponse{
			Message: handler.SOLUTION_REJECTED,
			Data:    submissionResult,
		})

		if err != nil {
			log.Println(err)
		}
	}

	room, err := h.isRoomExist(ws, payload.RoomId, h.Rooms)
	if err != nil {
		err := playerReq.Connection.WriteJSON(&entity.WSResponse{
			Message: handler.ROOM_IS_NOT_EXIST,
		})

		if err != nil {
			log.Println(err)
		}
	}

	h.broadcastMessageWhenSomeoneFinish(playerReq.Connection, room)

	h.insertDataToRedisLeaderboard(payload.RoomId, playerReq.PlayerId)

	h.Rooms[payload.RoomId] = h.insertDataLeaderboardToRoom(room)
	return nil
}

func (h *LeaderboardHandlerImpl) insertDataLeaderboardToRoom(room *entity.GameRoom) *entity.GameRoom {
	redisResult, err := h.RedisConnection.ZRevRangeWithScores(context.Background(), room.RoomId, 0, -1).Result()
	if err != nil {
		log.Println(err)
	}

	for i, v := range redisResult {
		room.Leaderboard = append(room.Leaderboard, &entity.PlayerLeaderboard{
			PlayerId:    v.Member.(string),
			DisplayName: h.findPlayerWithId(v.Member.(string), room.Players).DisplayName,
			Score:       v.Score,
			Position:    i + 1,
		})
	}

	return room
}

func (h *LeaderboardHandlerImpl) insertDataToRedisLeaderboard(roomId, playerId string) {
	err := h.RedisConnection.ZIncrBy(context.Background(), roomId, 100, playerId).Err()
	if err != nil {
		log.Println(err)
	}
}

func (h *LeaderboardHandlerImpl) isAllTestCasesPassed(playerConn *websocket.Conn, submissionResult *join_model.SubmissionTestResult) bool {
	for _, testResult := range submissionResult.TestResult {
		if testResult.Result != "PASS" {
			return false
		}
	}

	return true
}

func (h *LeaderboardHandlerImpl) getSubmissionResult(submissionId string, isCustomTest bool) *join_model.SubmissionTestResult {
	submissionResult := new(join_model.SubmissionTestResult)

	for {
		res, err := h.SubmissionUseCase.GetSubmission(submissionId, isCustomTest)
		if err != nil {
			log.Println(err)
			time.Sleep(70 * time.Millisecond)
			continue
		}

		if res.SubmissionStatus == "FINISHED" {
			submissionResult = res
			break
		}

		time.Sleep(70 * time.Millisecond)
	}

	return submissionResult
}

func (h *LeaderboardHandlerImpl) GetLeaderboard(c echo.Context) error {
	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ws.Close()

	roomId := c.QueryParam("room-id")
	playerId := c.QueryParam("player-id")
	room, err := h.isRoomExist(ws, roomId, h.Rooms)
	if err != nil {
		h.writeErrorResponse(ws, handler.ROOM_NOT_FOUND, nil)
		return err
	}

	player := h.findPlayerWithId(playerId, room.Players)
	h.writeSuccessResponse(player.Connection, handler.LEADERBOARD, room.Leaderboard)

	return nil
}

func (h *LeaderboardHandlerImpl) broadcastMessageWhenSomeoneFinish(playerReqConn *websocket.Conn, room *entity.GameRoom) {
	for _, player := range room.Players {
		if player.Connection != playerReqConn {
			err := player.Connection.WriteJSON(&entity.WSResponse{
				Message: handler.SOMEONE_FINISH,
				Data:    player.DisplayName,
			})

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (h *LeaderboardHandlerImpl) broadcastMessageWhenGameIsStart(room *entity.GameRoom) {
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
}

func (h *LeaderboardHandlerImpl) writeSuccessResponse(ws *websocket.Conn, message string, body any) {
	err := ws.WriteJSON(&entity.WSResponse{
		Message: message,
		Data:    body,
	})

	if err != nil {
		log.Println(err)
	}
}

func (h *LeaderboardHandlerImpl) writeErrorResponse(ws *websocket.Conn, message string, body any) {
	err := ws.WriteJSON(&entity.WSResponse{
		Message: message,
		Data:    body,
	})

	if err != nil {
		log.Println(err)
	}
}

func (h *LeaderboardHandlerImpl) broadcastMessageWhenRoomReadyToStart(room *entity.GameRoom) {
	if len(room.Players) == room.MaxUser {
		for _, player := range room.Players {
			err := player.Connection.WriteJSON(&entity.WSResponse{
				Message: handler.ROOM_IS_FULL_READY_TO_START,
				Data:    room,
			})

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (h *LeaderboardHandlerImpl) broadcastMessageWhenSomeoneJoinRoom(ws *websocket.Conn, room *entity.GameRoom) {
	for _, player := range room.Players {
		if player.Connection == ws {
			err := player.Connection.WriteJSON(&entity.WSResponse{
				Message: handler.SUCCESS_JOIN_ROOM,
				Data:    room,
			})

			if err != nil {
				log.Println(err)
			}

		} else {
			err := player.Connection.WriteJSON(&entity.WSResponse{
				Message: handler.SOMEONE_JOIN_ROOM,
				Data:    room,
			})

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (h *LeaderboardHandlerImpl) addNewPlayer(room *entity.GameRoom, newPlayer *entity.Player) *entity.GameRoom {
	room.Players = append(room.Players, newPlayer)
	return room
}

func (h *LeaderboardHandlerImpl) readJSON(ws *websocket.Conn, binder any) error {
	err := ws.ReadJSON(binder)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (h *LeaderboardHandlerImpl) isRoomAlreadyStart(ws *websocket.Conn, room *entity.GameRoom) bool {
	if room.IsGameStart {
		err := ws.WriteJSON(&entity.WSResponse{
			Message: handler.GAME_START,
		})

		if err != nil {
			log.Println(err)
		}

		return true
	}

	return false
}

func (h *LeaderboardHandlerImpl) isRoomFull(ws *websocket.Conn, room *entity.GameRoom) bool {
	if len(room.Players) == room.MaxUser {
		err := ws.WriteJSON(&entity.WSResponse{
			Message: handler.ROOM_IS_FULL,
		})

		if err != nil {
			log.Println(err)
		}

		return true
	}

	return false
}

func (h *LeaderboardHandlerImpl) isRoomExist(ws *websocket.Conn, roomId string, rooms map[string]*entity.GameRoom) (*entity.GameRoom, error) {
	room, found := rooms[roomId]
	if !found {
		err := ws.WriteJSON(&entity.WSResponse{
			Message: handler.ROOM_NOT_FOUND,
		})

		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return room, nil
}

func (h *LeaderboardHandlerImpl) findPlayerWithId(playerId string, players []*entity.Player) *entity.Player {
	for _, player := range players {
		if player.PlayerId == playerId {
			return player
		}
	}

	return nil
}
