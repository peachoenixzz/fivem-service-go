package playerlogs

import (
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
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
	ID      primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
	Event   string             `json:"event" bson:"event"`
	Content string             `json:"content" bson:"content"`
	Source  int                `json:"source" bson:"source"`
	Color   string             `json:"color" bson:"color"`
	Options struct {
		Public    bool `json:"public" bson:"public"`
		Important bool `json:"important" bson:"important"`
	} `json:"options" bson:"options"`
	Image     string    `json:"image" bson:"image"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Player    struct {
		Name        string `json:"name" bson:"name"`
		Identifiers struct {
			Ip       string `json:"ip" bson:"ip"`
			Steam    string `json:"steam" bson:"steam"`
			Discord  string `json:"discord" bson:"discord"`
			License  string `json:"license" bson:"license"`
			License2 string `json:"license2" bson:"license2"`
		} `json:"identifiers" bson:"identifiers"`
		Steam struct {
			Id     int    `json:"id" bson:"id"`
			Avatar string `json:"avatar" bson:"avatar"`
			Url    string `json:"url" bson:"url"`
		} `json:"steam" bson:"steam"`
	} `json:"player" bson:"player"`
	Hardware []string `json:"hardware" bson:"hardware"`
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
	logger.Info("get request event endpoint successfully")
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
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var mes Message
	mes, err = h.InsertMLog(req)
	logger.Info("prepare data to create successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, mes)
	}

	logger.Info("create successfully")
	return c.JSON(http.StatusCreated, mes)
}

func (h Handler) CaseEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := c.Param("steamid")
	event := c.Param("event")
	logger.Info("prepare log")
	res, err := h.LogCaseEventAndSteamID(steamID, event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event case and steam id endpoint")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AllEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := c.Param("steamid")
	res, err := h.LogAllEventAndSteamID(steamID)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event and steamid endpoint")
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

	logger.Info("get event endpoint successfully")
	return c.JSON(http.StatusOK, res)
}
