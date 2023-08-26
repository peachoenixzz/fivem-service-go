package playerquest

import (
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Message struct {
	Message interface{}
}

type ResponseValidateItem struct {
	LimitType      string `json:"limit_type"`
	Name           string `json:"item_name"`
	Category       string `json:"category"`
	MaxLimit       int64  `json:"max_limit"`
	Point          int64  `json:"point"`
	RemainQuantity int64  `json:"remaining_quantity"`
	ExpireDateItem int    `json:"expire_date_item"`
}

type ResponseRequireQuestPlayer struct {
	Require     int64 `json:"require"`
	WeightLevel int64 `json:"weightlevel"`
}

// func (h Handler) CreateQuestPlayer(c echo.Context) error {
//
//		return c.JSON(http.StatusOK, res)
//	}
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
