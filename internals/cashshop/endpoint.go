package cashshop

import (
	"context"
	"database/sql"
	"fmt"
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

type ResponseInitCashShop struct {
	Identifier    string `json:"identifier"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Point         string `json:"point"`
	ExpireDateVip string `json:"expire_date_vip"`
}

type RequestBuyItem struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

type ResponseItemCashShop struct {
	LimitType      string `json:"limit_type"`
	Name           string `json:"item_name"`
	LabelName      string `json:"label_name"`
	Description    string `json:"description"`
	PromotionFlag  string `json:"promotion_flag"`
	MaxLimit       int64  `json:"max_limit"`
	Point          int64  `json:"point"`
	RemainQuantity int64  `json:"remaining_quantity"`
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
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
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

	if req.Quantity <= 0 {
		tx.Rollback()
		logger.Error("Invalid quantity requested.", zap.Int64("Requested", req.Quantity))
		return c.JSON(http.StatusOK, Message{Message: "fail"})
	}

	res, err := h.ValidatePurchaseItem(tx, req, playerInfo.Identifier)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	if req.Quantity > res.RemainQuantity && res.RemainQuantity != -1 { // Assuming `req.Quantity` is the user input
		tx.Rollback()
		logger.Error("User requested more than available quantity.", zap.Int64("Requested", req.Quantity), zap.Int64("Available", res.RemainQuantity))
		return c.JSON(http.StatusOK, Message{Message: "fail"})
	}

	logger.Info(fmt.Sprintf("Name : %v Point : %v Limit : %v Type : %v Remain : %v Expire : %v Category : %v", res.Name, res.Point, res.MaxLimit, res.LimitType, res.RemainQuantity, res.ExpireDateItem, res.Category))
	if HandleLimitType(res) {
		logger.Info("prepare PurchaseItem")
		count, err := h.PurchaseItem(tx, req, playerInfo.Identifier)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to Update record:", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
		}

		if count > 0 {
			switch res.Category {
			case "normal":
				err = buyGeneralItems(h, tx, req, res, playerInfo)
				if err != nil {
					return c.JSON(http.StatusOK, Message{Message: "fail"})
				}
				tx.Commit()
				return c.JSON(http.StatusOK, "success")
			case "vehicle":
				err = buyVehicleItems(h, tx, req, res, playerInfo)
				if err != nil {
					return c.JSON(http.StatusOK, Message{Message: "fail"})
				}
				tx.Commit()
				return c.JSON(http.StatusOK, "success")
			}
		}
	}
	tx.Rollback()
	return c.JSON(http.StatusOK, Message{Message: "fail"})
}

func buyVehicleItems(h Handler, tx *sql.Tx, req RequestBuyItem, res ResponseValidateItem, playerInfo *mw.JwtCustomClaims) error {
	logger := mlog.Logg
	if HandleDateExpire(res.ExpireDateItem) {
		if req.Quantity > 1 {
			tx.Rollback()
			logger.Error("User requested more than available quantity.", zap.Int64("Requested", req.Quantity), zap.Int64("Available", res.RemainQuantity))
			return fmt.Errorf("user requested more than available quantity")
		}
		err := h.InsertExpireDateVehicle(tx, res, playerInfo.Identifier)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to Update record:", zap.Error(err))
			return fmt.Errorf("failed to update record")
		}
		for i := 0; i < int(req.Quantity); i++ {
			err := h.InsertHistoryPurchaseItem(tx, res, playerInfo.Identifier)
			if err != nil {
				tx.Rollback()
				logger.Error("Failed to Update record:", zap.Error(err))
				return fmt.Errorf("failed to update record")
			}
		}
	}
	return nil
}

func buyGeneralItems(h Handler, tx *sql.Tx, req RequestBuyItem, res ResponseValidateItem, playerInfo *mw.JwtCustomClaims) error {
	logger := mlog.Logg
	if HandleDateExpire(res.ExpireDateItem) {
		if req.Quantity > 1 {
			tx.Rollback()
			logger.Error("User requested more than available quantity.", zap.Int64("Requested", req.Quantity), zap.Int64("Available", res.RemainQuantity))
			return fmt.Errorf("user requested more than available quantity")
		}
		err := h.InsertExpireDateItem(tx, res, playerInfo.Identifier)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to Update record:", zap.Error(err))
			return fmt.Errorf("failed to update record")
		}
	}
	for i := 0; i < int(req.Quantity); i++ {
		err := h.InsertHistoryPurchaseItem(tx, res, playerInfo.Identifier)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to Update record:", zap.Error(err))
			return fmt.Errorf("failed to update record")
		}
	}
	err := h.InsertGivePlayerItem(tx, res, req, playerInfo.Identifier)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return fmt.Errorf("failed to update record")
	}
	return nil
}
