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
	"os"
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
	Quantity int
}

type GachaponItem struct {
	Item       Item
	GachaponID int
	PullRate   float64
}

type ResponseItemInGashapon struct {
	Name      string `json:"name"`
	LabelName string `json:"label_name"`
}

type RequestGashaponName struct {
	Name string `json:"name"`
}

type ResponseGiveItemStatus struct {
	InSlot int `json:"in_slot"`
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

func (h Handler) GetInSlotGiveItemsInGachaponEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	req := RequestGashaponName{}
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
	//user := c.Get("user").(*jwt.Token)
	//playerInfo := user.Claims.(*mw.JwtCustomClaims)
	req := RequestGashaponName{}
	err := c.Bind(&req)
	gci, err := h.GetGashaponItemsRate(context.Background(), req)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}

	// Run the gachapon draw 300 times
	itemStats := make(map[string]map[string]float64)

	// Run the gachapon draw 300 times
	for i := 0; i < 500; i++ {
		drawnItem, pullRate := handleRandGachaponItems(gci)
		if drawnItem != nil {
			if _, exists := itemStats[drawnItem.Name]; !exists {
				itemStats[drawnItem.Name] = map[string]float64{
					"count":    0,
					"pullRate": pullRate,
				}
			}
			itemStats[drawnItem.Name]["count"]++
		}
	}

	// Create and open the text file
	file, err := os.Create(fmt.Sprintf("gachapon_summary_%v.txt", req.Name))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Print and write the summary to the text file
	fmt.Fprintln(file, "Gachapon Summary:")
	fmt.Println("Gachapon Summary:")
	for itemName, stats := range itemStats {
		actualRate := (stats["count"] / 300.0) * 100
		summary := fmt.Sprintf("%s: Count %.0f, Expected PullRate %f%%, Actual PullRate %f%%",
			itemName, stats["count"], stats["pullRate"]*100, actualRate)
		fmt.Fprintln(file, summary)
		fmt.Println(summary)
	}

	return nil
}
