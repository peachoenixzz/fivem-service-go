package gachapon

import (
	"fmt"
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

type Item struct {
	Name     string
	Category string
	ItemId   string
	Quantity int
}

type ItemInsert struct {
	Name       string
	Category   string
	ItemId     string
	Quantity   int
	GachaponID int
}

type GachaponItem struct {
	Item       Item
	GachaponID int
	PullRate   float64
}

type ResponseItemInGachapon struct {
	Name      string `json:"name"`
	LabelName string `json:"label_name"`
	ItemId    string `json:"id"`
}

type RequestGachaponName struct {
	Name string `json:"name"`
}

type RequestOpenGachapon struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type ResponseGiveItemStatus struct {
	InSlot int `json:"in_slot"`
}

const (
	maxQuantity = 500
)

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
	req := RequestGachaponName{}
	err := c.Bind(&req)
	ig, err := h.GetItemsInGachapon(context.Background(), req)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	return c.JSON(http.StatusOK, ig)
}

func (h Handler) GetInSlotGiveItemsInGachaponEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	req := RequestGachaponName{}
	err := c.Bind(&req)
	st, err := h.GetInSlotGiveItemsGachapon(context.Background(), req, playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	return c.JSON(http.StatusOK, st)
}

func (h Handler) OpenGachaponEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	req := RequestOpenGachapon{}
	err := c.Bind(&req)

	if req.Quantity > maxQuantity {
		logger.Error(fmt.Sprintf("wrong quantity to request (maxPull reason) (%v)", playerInfo.Identifier))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("wrong quantity to request (%v)", playerInfo.Identifier))
	}

	if req.Quantity <= 0 {
		logger.Error(fmt.Sprintf("wrong quantity to request (%v)", playerInfo.Identifier))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("wrong quantity to request (%v)", playerInfo.Identifier))
	}

	pi, err := h.QueryPlayerItem(context.Background(), playerInfo.Identifier)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}

	if pi[req.Name] < req.Quantity {
		logger.Error(fmt.Sprintf("player have items less than request (%v) (actual %v , expected %v)", playerInfo.Identifier, pi[req.Name], req.Quantity))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("player have items less than request (%v)", playerInfo.Identifier))
	}

	gci, err := h.GetGashaponItemsRate(context.Background(), req)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}

	itemRand := make(map[string]map[string]any)
	for i := 0; i < req.Quantity; i++ {
		drawnItem, gid, cg := handleRandGachaponItems(gci)
		if drawnItem != nil {
			if _, exists := itemRand[drawnItem.ItemId]; !exists {
				itemRand[drawnItem.ItemId] = map[string]any{
					"count":       0,
					"gachapon_id": gid,
					"category":    cg,
				}
			}
			itemRand[drawnItem.ItemId]["count"] = itemRand[drawnItem.ItemId]["count"].(int) + 1
		}
	}

	itemsInsert, items := handleRandResponseAndInsertGachapon(itemRand, gci)
	tx, err := h.MysqlDB.Begin()
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if err != nil {
		logger.Error("Failed transaction:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	err = h.InsertItemPrepareGivePlayer(tx, itemsInsert, req, playerInfo.Identifier)
	if err != nil {
		logger.Error("Failed transaction insert give item:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error")
	}

	return c.JSON(http.StatusOK, items)
}
