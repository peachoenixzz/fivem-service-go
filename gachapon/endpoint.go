package gachapon

import (
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
)

type Message struct {
	Message interface{}
}

type ResponsePlayerGachapon struct {
	Name      string `json:"name"`
	LabelName string `json:"label_name"`
	Quantity  int    `json:"quantity"`
}

type AllGachapon struct {
	Name      string `json:"name"`
	LabelName string `json:"label_name"`
}

type ResponseItemInGashapon struct {
	Name      string `json:"name"`
	LabelName string `json:"label_name"`
}

type RequestGashaponName struct {
	Name string `json:"name"`
}

func (h Handler) GetPlayerGachaponEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	pi, err := h.QueryPlayerItem(context.Background(), playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}

	ag, err := h.GetAllGachapon(context.Background())
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}

	res := handleGachaponPlayer(pi, ag)
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) GetItemsInGachaponEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	req := RequestGashaponName{}
	err := c.Bind(&req)
	ig, err := h.GetItemsInGachapon(context.Background(), req)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	return c.JSON(http.StatusOK, ig)
}
