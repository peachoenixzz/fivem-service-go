package playerquest

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

type ResponseRequireQuestPlayer struct {
	Quantity    int64 `json:"quantity"`
	WeightLevel int64 `json:"weightlevel"`
}

type PlayerItems map[string]interface{}

type ResponseQuestItem struct {
	Name string
	Rare string
}

type ResponseItemComparison struct {
	ItemName   string `json:"item_name"`
	Comparison string `json:"comparison"`
}

type ResponseSelectedItem struct {
	Name     string
	Rare     string
	Quantity int64
}

type ResponsePlayerQuestItem struct {
	Name     string
	Quantity int64
}

func (h Handler) GetRequireQuestPlayer(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	res, err := h.QueryRequireQuest(c, playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) CreateQuestPlayer(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	rqi, err := h.QueryQuestItem(context.Background())
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}

	if h.GetStateQuest(context.Background(), playerInfo.Identifier) {
		rsi := handleQuestItem(rqi)
		h.InsertSelectQuestItem(rsi, playerInfo.Identifier)
		logger.Info("player get quest success")
		return c.JSON(http.StatusOK, "success")
	}

	logger.Info("player already get quest")
	return c.JSON(http.StatusOK, "failed")
}

func (h Handler) GetComparePlayerItemAndQuestItem(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	pi, err := h.QueryPlayerItem(context.Background(), playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	rpqi, err := h.GetPlayerQuestItem(context.Background(), playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	handleComparePlayerAndQuestItem(pi, rpqi)
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, "eiei")
}
