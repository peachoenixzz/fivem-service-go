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
	WeightLevel int64 `json:"weight_level"`
	CardAItem   int   `json:"card_a_item"`
}

type ResponseQuestItem struct {
	Name string
	Rare string
}

type ResponseItemComparison struct {
	ItemName             string `json:"item_name"`
	LabelName            string `json:"label_name"`
	Comparison           string `json:"comparison"`
	PlayerItemQuantity   int    `json:"player_item_quantity"`
	QuestRequireQuantity int64  `json:"quest_require_quantity"`
}

type ResponseSelectedItem struct {
	Name     string
	Rare     string
	Quantity int64
}

type ResponsePlayerQuestItem struct {
	ItemName  string
	LabelName string
	Quantity  int64
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
	pi, err := h.QueryPlayerItem(context.Background(), playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	res = handleCardAItem(res, pi)
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) GetStatusQuest(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	if h.GetStateQuest(context.Background(), playerInfo.Identifier) {
		logger.Info("player not ready get quest ")
		return c.JSON(http.StatusOK, Message{Message: "not_ready_quest"})
	}

	logger.Info("player already get quest")
	return c.JSON(http.StatusOK, Message{Message: "already_quest"})
}

func (h Handler) ResetQuestPlayer(c echo.Context) error {
	logger := mlog.Logg
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)

	tx, err := h.MysqlDB.Begin()
	if err != nil {
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	err = h.ResetQuest(tx, playerInfo.Identifier)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Database Err:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	return c.JSON(http.StatusOK, Message{Message: "Reset Successfully"})

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
		return c.JSON(http.StatusOK, Message{Message: "success"})
	}

	logger.Info("player already get quest")
	return c.JSON(http.StatusOK, Message{Message: "failed"})
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
	res := handleComparePlayerAndQuestItem(pi, rpqi)
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}
