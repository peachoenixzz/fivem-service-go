package cashshop

import (
	//"fmt"
	//"github.com/golang-jwt/jwt/v4"
	//mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Message struct {
	Status  int
	Message interface{}
}

type ResponseInitCashShop struct {
	Identifier string `json:"identifier"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Point      string `json:"point"`
}

type Request struct {
}

type RequestUpdatePoint struct {
	Identifier string `json:"identifier"`
	DiscordID  string `json:"discord_id"`
	CashPoint  int64  `json:"cashPoint"`
}

func (h Handler) GetInitCashShopEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	discordID := c.Param("discordID")
	//user := c.Get("user").(*jwt.Token)
	//playerInfo := user.Claims.(*mw.JwtCustomClaims)
	//fmt.Println("JOB", playerInfo.Job)
	//fmt.Println("identifier", playerInfo.Identifier)
	//fmt.Println("group", playerInfo.Group)
	res, err := h.getInitCashShop(c, discordID)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) UpdateCashPointEndPoint(c echo.Context) error {
	logger := mlog.Logg
	discordID := c.Param("discordID")
	var req RequestUpdatePoint
	err := c.Bind(&req)
	if err != nil {
		logger.Error("Failed to bind request:", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to bind request")
	}

	tx, err := h.MysqlDB.Begin()
	if err != nil {
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	err = h.UpdateCashPoint(tx, req, discordID)
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

	return c.JSON(http.StatusOK, Message{Message: "Update Cash Point Successfully"})

}
