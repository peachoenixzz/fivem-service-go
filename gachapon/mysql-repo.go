package gachapon

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}

func (h Handler) GetAllGachapon(ctx context.Context) ([]AllGachapon, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	stmtStr := `SELECT i.label, tg.name  FROM TB_GACHAPON tg 
				INNER JOIN items i 
				WHERE tg.name = i.name ;
	`

	// Create a prepared statement
	logger.Info("mysql prepare query TB_GACHAPON")

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []AllGachapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []AllGachapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	var items []AllGachapon
	for rows.Next() {
		var item AllGachapon
		err := rows.Scan(&item.LabelName, &item.Name)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []AllGachapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, item)
	}
	return items, nil
}

func (h Handler) GetItemsInGachapon(ctx context.Context, req RequestGashaponName) ([]ResponseItemInGashapon, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query item in gachapon name")
	stmtStr := `SELECT tgi.name,CONCAT(i.label," (จำนวน ",tgi.quantity ,")")  FROM TB_GACHAPON tg 
				INNER JOIN TB_GACHAPON_ITEMS tgi 
				ON tgi.gachapon_id = tg.gachapon_id 
				INNER JOIN items i 
				ON i.name = tgi.name 
				WHERE
				tg.name  = ?
	`

	// Create a prepared statement
	logger.Info("mysql prepare query TB_GACHAPON")

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemInGashapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer stmt.Close()

	args := []interface{}{
		req.Name,
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemInGashapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	var items []ResponseItemInGashapon
	for rows.Next() {
		var item ResponseItemInGashapon
		err := rows.Scan(&item.Name, &item.LabelName)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []ResponseItemInGashapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, item)
	}
	return items, nil
}

func (h Handler) QueryPlayerItem(ctx context.Context, discordID string) (map[string]int, error) {
	logger := mlog.Logg
	stmtStr := "SELECT inventory FROM users u WHERE u.identifier = ?"
	logger.Info("mysql prepare query Discord ID")
	var playerItems map[string]int
	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("sql error", zap.Error(err))
		return playerItems, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	logger.Info("query Row Discord Id")
	var itemStr string
	err = stmt.QueryRow(discordID).Scan(&itemStr)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return playerItems, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("after query row and ready to return Data")

	//bad idea to handle [] from ESX Framework when no items
	if itemStr == `[]` {
		itemStr = `{"mockup":1}`
	}
	if err := json.Unmarshal([]byte(itemStr), &playerItems); err != nil {
		logger.Error("Failed to parse JSON item data:", zap.Error(err))
		return playerItems, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	return playerItems, nil
}

func (h Handler) GetInSlotGiveItemsGachapon(ctx context.Context, req RequestGashaponName, discordID string) (ResponseGiveItemStatus, error) {
	logger := mlog.Logg
	query := `
		SELECT count(1) FROM TB_GACHAPON tg 
		INNER JOIN TB_GIVE_ITEMS_GACHAPON tgig 
		ON tg.gachapon_id  = tgig.gachapon_id 
		WHERE tg.name = ?
		AND tgig.status = 'pending'
		AND discord_id = ?;
`

	// Create a prepared statement
	logger.Info("mysql prepare query status gachapon")
	stmt, err := h.MysqlDB.PrepareContext(ctx, query)
	if err != nil {
		logger.Error("sql error", zap.Error(err))
	}
	defer stmt.Close()

	args := []interface{}{
		req.Name,
		discordID,
	}

	logger.Info("query status gachapon player")
	var gis ResponseGiveItemStatus
	err = stmt.QueryRow(args...).Scan(&gis.InSlot)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return ResponseGiveItemStatus{}, err
	}
	logger.Info("after query row and ready to return Data")
	return gis, nil
}
