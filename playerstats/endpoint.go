package playerstats

import (
	"context"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Message struct {
	Message interface{} `json:"message"`
}

type ResponseVipItems struct {
	ItemName  string `json:"item_name"`
	ItemLabel string `json:"item_label"`
}

func (h Handler) GetAllMoney(c echo.Context) error {
	logger := mlog.L(c)
	err := h.AllMoney(context.Background())
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("get request event endpoint successfully")

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, "success")
}

func (h Handler) GetVehicleByModel(c echo.Context) error {
	logger := mlog.L(c)
	err := h.VehicleByModel(context.Background())
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("get request event endpoint successfully")

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, "success")
}
