package playerlogin

import (
	"fmt"
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
	logger.Info("prepare log")
	fmt.Printf("steam id is %s", req.Identifier)

	res, err := h.PlayerIdentify(c.Request().Context(), req)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("GetPlayerIdentify endpoint end")
	return mw.LoginSuccess(c, mw.Response(res))
}
