package playerlogs

import (
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (h Handler) CustomLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	req := RequestCustomLog{}
	err := c.Bind(&req)
	logger.Info("get request custom log endpoint successfully")
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var res []Response
	res, err = h.CustomMLog(req)
	logger.Info("prepare data to setup successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	logger.Info("get custom fivem log successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) GetFiveMLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	res, err := h.FiveMLog()
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AddFiveMLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	req := RequestInsert{}
	err := c.Bind(&req)
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

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
