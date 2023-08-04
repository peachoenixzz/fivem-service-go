package playeridentifier

import (
	"fmt"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Message struct {
	Message interface{} `json:"message"`
}

type Item struct {
	ItemName   string `json:"item_name"`
	Quantity   int    `json:"quantity"`
	ExpireDate int    `json:"expire_date"`
}

type RequestUpdateVip struct {
	DiscordID      string `json:"discord_id"`
	Priority       string `json:"priority"`
	Identifier     string `json:"identifier"`
	VipPoint       int64  `json:"vip_point"`
	ExtraPoint     int64  `json:"extra_point"`
	PermanentPoint int64  `json:"permanent_point"`
	ExpireItems    []Item `json:"expire_items"`
}

type Response struct {
	DiscordID      string `json:"discord_id"`
	SteamID        string `json:"steam_id"`
	Priority       string `json:"priority"`
	ExpireDate     string `json:"expire_date"`
	LastUpdated    string `json:"last_updated"`
	Identifier     string `json:"identifier"`
	VipPoint       int64  `json:"vip_point"`
	ExtraPoint     int64  `json:"extra_point"`
	PermanentPoint int64  `json:"permanent_point"`
}

func (h Handler) UpdateVIPPoint(c echo.Context) error {
	logger := mlog.L(c)
	req := RequestUpdateVip{}
	err := c.Bind(&req)
	if err != nil {
		logger.Error("bad Request Update Vip body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad RequestUpdateVip body", err.Error())
	}
	logger.Info("get RequestUpdateVip event endpoint successfully")
	discordID := c.Param("discordID")
	tx, err := h.MysqlDB.Begin()
	err = h.UpdateVipPointByPlayerDiscord(tx, req, discordID)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	for _, item := range req.ExpireItems {
		logger.Info("prepare to Handle Vip player item")
		if HandleDateExpire(item.ExpireDate) {
			logger.Info("item have expire date")
			err = h.InsertExpireDateItem(tx, item, req.DiscordID)
			if err != nil {
				tx.Rollback()
				logger.Error("Failed to Insert record:", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
			}
		}
		err = h.InsertPlayerItems(tx, item, req.DiscordID)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to Insert record:", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		fmt.Println("Item Name : ", item.ItemName)
		fmt.Println("Quantity : ", item.Quantity)
		fmt.Println("Exp Date : ", item.ExpireDate)
	}
	err = tx.Commit()
	if err != nil {
		logger.Error("Database Err : ", zap.Error(err))
		tx.Rollback()
		return echo.NewHTTPError(http.StatusBadRequest, "Database Error : ", err.Error())
	}
	return c.JSON(http.StatusOK, Message{Message: "Update VIP and insert items Successfully"})
}

func (h Handler) GetPlayerDiscordID(c echo.Context) error {
	logger := mlog.L(c)
	discordID := c.Param("discordID")
	res, err := h.QueryPlayerDiscord(c, discordID)
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}
