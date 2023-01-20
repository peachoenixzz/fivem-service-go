package fivemlogs

import (
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Message struct {
	Status  int
	Message interface{}
}

type Response struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Event   string             `bson:"event"`
	Content string             `bson:"content"`
	Source  int                `bson:"source"`
	Color   string             `bson:"color"`
	Options struct {
		Public    bool `bson:"public"`
		Important bool `bson:"important"`
	} `json:"options"`
	Image     string    `bson:"image"`
	Timestamp time.Time `bson:"timestamp"`
	Player    struct {
		Name        string `bson:"name"`
		Identifiers struct {
			Ip       string `bson:"ip"`
			Steam    string `bson:"steam"`
			Discord  string `bson:"discord"`
			License  string `bson:"license"`
			License2 string `bson:"license2"`
		} `bson:"identifiers"`
		Steam struct {
			Id     int    `bson:"id"`
			Avatar string `bson:"avatar"`
			Url    string `bson:"url"`
		} `bson:"steam"`
	} `bson:"player"`
	Hardware []string `bson:"hardware"`
}

type Request struct {
	Event   string `json:"event"`
	Content string `json:"content"`
	Source  int    `json:"source"`
	Color   string `json:"color"`
	Options struct {
		Public    bool `json:"public"`
		Important bool `json:"important"`
	} `json:"options"`
	Image     string    `json:"image"`
	Timestamp time.Time `json:"timestamp"`
	Player    struct {
		Name        string `json:"name"`
		Identifiers struct {
			Ip       string `json:"ip"`
			Steam    string `json:"steam"`
			Discord  string `json:"discord"`
			License  string `json:"license"`
			License2 string `json:"license2"`
		} `json:"identifiers"`
		Steam struct {
			Id     int    `json:"id"`
			Avatar string `json:"avatar"`
			Url    string `json:"url"`
		} `json:"steam"`
	} `json:"player"`
	Hardware []string `json:"hardware"`
}

func (h Handler) GetFiveMLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	res, err := h.FiveMLog()
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AddFiveMLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	req := Request{}
	err := c.Bind(&req)
	//req.Event = strings.ToLower(req.Event)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var mes Message
	mes, err = h.InsertMLog(req)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, mes)
	}

	logger.Info("create successfully")
	return c.JSON(http.StatusCreated, mes)
}

func (h Handler) CaseEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := cast.ToInt(c.Param("steamid"))
	event := c.Param("event")
	res, err := h.LogByEventAndSteamID(steamID, event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("Get Log by steamid and event successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AllEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := cast.ToInt(c.Param("steamid"))
	res, err := h.LogAllEventAndSteamID(steamID)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("Get Log by steamid and event successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) ByEventEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	event := c.Param("event")
	res, err := h.LogCaseEventAll(event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("Get Log by steamid and event successfully")
	return c.JSON(http.StatusOK, res)
}
