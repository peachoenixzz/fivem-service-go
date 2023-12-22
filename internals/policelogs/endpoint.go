package policelogs

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Message struct {
	Status  int
	Message interface{}
}

type Response struct {
	ArrestPlayerName      string        `json:"arrest_player_name"`
	PolicePlayerName      string        `json:"police_player_name"`
	ArrestSteamPlayerName string        `json:"arrest_steam_player_name"`
	PoliceSteamPlayerName string        `json:"police_steam_player_name"`
	ArrestJobPlayer       string        `json:"arrest_job_player"`
	ArrestSexPlayer       string        `json:"arrest_sex_player"`
	Case                  []interface{} `json:"case"`
	CaseQuantity          []interface{} `json:"case_quantity"`
	CaseCustom            []interface{} `json:"case_custom"`
	TimeCustom            []interface{} `json:"time_custom"`
	FineCustom            []interface{} `json:"fine_custom"`
	AllMiliSec            int64         `json:"all_milisec"`
	AllMinute             int64         `json:"all_mins"`
	AllFine               int64         `json:"all_fine"`
	PoliceDecreaseTime    int64         `json:"police_decrease_time"`
}

type Request struct {
	ArrestPlayerName      string        `json:"arrest_player_name"`
	PolicePlayerName      string        `json:"police_player_name"`
	ArrestSteamPlayerName string        `json:"arrest_steam_player_name"`
	PoliceSteamPlayerName string        `json:"police_steam_player_name"`
	ArrestJobPlayer       string        `json:"arrest_job_player"`
	ArrestSexPlayer       string        `json:"arrest_sex_player"`
	Case                  []interface{} `json:"case"`
	CaseQuantity          []interface{} `json:"case_quantity"`
	CaseCustom            []interface{} `json:"case_custom"`
	TimeCustom            []interface{} `json:"time_custom"`
	FineCustom            []interface{} `json:"fine_custom"`
	AllMiliSec            int64         `json:"all_milisec"`
	AllMinute             int64         `json:"all_mins"`
	AllFine               int64         `json:"all_fine"`
	PoliceDecreaseTime    int64         `json:"police_decrease_time"`
}

func (h Handler) GetPoliceLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	res, err := h.PoliceLog()
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AddPoliceLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	var req Request
	err := c.Bind(&req)
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	fmt.Println("JOB", playerInfo.Job)
	fmt.Println("identifier", playerInfo.Identifier)
	fmt.Println("group", playerInfo.Group)
	//if playerInfo.Job == "unemployed" {
	//	return c.String(http.StatusNonAuthoritativeInfo, "Wrong Access :( ")
	//}

	var mes Message
	mes, err = h.InsertMLog(req)
	logger.Info("prepare data to create successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, mes)
	}

	logger.Info("create successfully")
	return c.JSON(http.StatusCreated, mes)
}

func (h Handler) CaseEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := c.Param("steamid")
	event := c.Param("event")
	logger.Info("prepare log")
	res, err := h.LogCaseEventAndSteamID(steamID, event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event case and steam id endpoint")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AllEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := c.Param("steamid")
	res, err := h.LogAllEventAndSteamID(steamID)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event and steamid endpoint")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) ByEventEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	event := c.Param("event")
	res, err := h.LogCaseEventAll(event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event endpoint successfully")
	return c.JSON(http.StatusOK, res)
}
