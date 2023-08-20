package cashshop

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	//"fmt"
	//"github.com/golang-jwt/jwt/v4"
	//mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Message struct {
	Message interface{}
}

type ResponseInitCashShop struct {
	Identifier    string `json:"identifier"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Point         string `json:"point"`
	ExpireDateVip string `json:"expire_date_vip"`
}

type RequestBuyItem struct {
	Name string `json:"name"`
}

type ResponseItemCashShop struct {
	LimitType      string `json:"limit_type"`
	Name           string `json:"item_name"`
	MaxLimit       int64  `json:"max_limit"`
	Point          int64  `json:"point"`
	RemainQuantity int64  `json:"remaining_quantity"`
}

type ResponseValidateItem struct {
	LimitType      string `json:"limit_type"`
	Name           string `json:"item_name"`
	MaxLimit       int64  `json:"max_limit"`
	Point          int64  `json:"point"`
	RemainQuantity int64  `json:"remaining_quantity"`
}

func (h Handler) GetInitCashShopEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	res, err := h.getInitCashShop(c, playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) GetCashShopItemEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	res, err := h.GetCashShopItem(context.Background(), playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	logger.Info("get result item cash shop successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) BuyCashShopEndPoint(c echo.Context) error {
	logger := mlog.Logg
	var req RequestBuyItem
	err := c.Bind(&req)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	if err != nil {
		logger.Error("Failed to bind request:", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to bind request")
	}

	tx, err := h.MysqlDB.Begin()
	if err != nil {
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	res, err := h.ValidatePurchaseItem(tx, req, playerInfo.Identifier)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}
	logger.Info(fmt.Sprintf("Name : %v Point : %v Limit : %v Type : %v Remain : %v", res.Name, res.Point, res.MaxLimit, res.LimitType, res.RemainQuantity))
	if HandleLimitType(res) {
		logger.Info("prepare PurchaseItem")
		count, err := h.PurchaseItem(tx, req, playerInfo.Identifier)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to Update record:", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
		}
		if count > 0 {

		}

		logger.Info(fmt.Sprintf("update Count : %v", count))
		ms := HandleMessage(count)
		if ms.Message == "success" {
			logger.Info("prepare PurchaseItem Success")
			tx.Commit()
			return c.JSON(http.StatusOK, ms)
		}
		tx.Rollback()
		return c.JSON(http.StatusOK, ms)
	}
	tx.Rollback()
	return c.JSON(http.StatusOK, Message{Message: "fail"})
}
