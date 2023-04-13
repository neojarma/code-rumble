package leaderboard_handler

import "github.com/labstack/echo/v4"

type LeaderboardHandler interface {
	CreateRoom(c echo.Context) error
	JoinRoom(c echo.Context) error
	StartGame(c echo.Context) error
	SubmitCode(c echo.Context) error
	GetLeaderboard(c echo.Context) error
}
