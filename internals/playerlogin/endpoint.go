package playerlogin

import (
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
	Identifier string `json:"identifier"`
	Job        string `json:"job"`
	Group      string `json:"group"`
}

type Request struct {
	Identifier string `json:"identifier"`
}

func (h Handler) GetPlayerIdentify(c echo.Context) error {
	logger := mlog.L(c)
	var req Request
	logger.Info("prepare to bind request to struct request")
	if err := c.Bind(&req); err != nil {
		logger.Error("Bind Err: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	logger.Info("prepare to check discord ID", zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
	res, err := h.PlayerIdentify(c, req)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("GetPlayerIdentify endpoint end", zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
	return mw.LoginSuccess(c, mw.Response(res))
}
