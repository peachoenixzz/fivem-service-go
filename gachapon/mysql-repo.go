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

func (h Handler) GetItemsInGachapon(ctx context.Context, req RequestGachaponName) ([]ResponseItemInGachapon, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query item in gachapon name")
	stmtStr := `SELECT tgi.name,CONCAT(i.label," (จำนวน ",tgi.quantity ,")"),tgi.gachapon_item_id  FROM TB_GACHAPON tg 
				INNER JOIN TB_GACHAPON_ITEMS tgi 
				ON tgi.gachapon_id = tg.gachapon_id 
				INNER JOIN items i 
				ON i.name = tgi.name 
				WHERE
				tg.name  = ?
				ORDER BY pull_rate DESC
	`

	// Create a prepared statement
	logger.Info("mysql prepare query TB_GACHAPON")

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemInGachapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer stmt.Close()

	args := []interface{}{
		req.Name,
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemInGachapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	var items []ResponseItemInGachapon
	for rows.Next() {
		var item ResponseItemInGachapon
		err := rows.Scan(&item.Name, &item.LabelName, &item.ItemId)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []ResponseItemInGachapon{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, item)
	}
	return items, nil
}

func (h Handler) QueryPlayerItem(ctx context.Context, discordID string) (map[string]int, error) {
	logger := mlog.Logg
	stmtStr := "SELECT inventory FROM users u WHERE u.identifier = ?"
	logger.Info("mysql prepare query player item on inventory")
	var playerItems map[string]int
	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("sql error", zap.Error(err))
		return playerItems, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer stmt.Close()

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

func (h Handler) GetInSlotGiveItemsGachapon(ctx context.Context, req RequestGachaponName, discordID string) (ResponseGiveItemStatus, error) {
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

func (h Handler) InsertItemPrepareGivePlayer(tx *sql.Tx, i []ItemInsert, req RequestOpenGachapon, discordID string) error {
	logger := mlog.Logg
	logger.Info("prepare to do stmt history purchase item")
	for _, v := range i {
		stmtStr := `
		INSERT
		INTO TB_GIVE_ITEMS_GACHAPON (discord_id,item_name,quantity,status,gachapon_id,category,gachapon_name,created_date,last_update)
		VALUES (?,?,?,?,?,?,?,SYSDATE(),SYSDATE())
	`
		args := []interface{}{
			discordID,
			v.Name,
			v.Quantity,
			"pending",
			v.GachaponID,
			v.Category,
			req.Name,
		}

		//for _, v := range args {
		//	fmt.Println(v)
		//}

		logger.Info("prepare to Insert history purchase item")
		r, err := tx.Exec(stmtStr, args...)
		if err != nil {
			logger.Error("Failed to Insert gachapon Item record:", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}

		logger.Info("prepare to Insert history purchase item")
		rowsAffected, err := r.RowsAffected()
		if err != nil {
			logger.Error("Failed to Insert history purchase record: ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		logger.Info("Row Affected Update Insert history purchase table", zap.Int64("row affected", rowsAffected))
	}
	return nil
}

func (h Handler) GetGashaponItemsRate(ctx context.Context, req RequestOpenGachapon) ([]GachaponItem, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query item in gachapon name")
	stmtStr := `SELECT  tg.gachapon_id,i.name, tgi.pull_rate , tgi.quantity,tgi.category,tgi.gachapon_item_id
		FROM TB_GACHAPON_ITEMS tgi  
		INNER JOIN items i 
		ON tgi.name = i.name 
		INNER JOIN TB_GACHAPON tg 
		ON tg.gachapon_id = tgi.gachapon_id 
		WHERE tg.name = ?;
	`

	// Create a prepared statement
	logger.Info("mysql prepare query TB_GACHAPON")

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []GachaponItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer stmt.Close()

	args := []interface{}{
		req.Name,
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []GachaponItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	var gci []GachaponItem
	for rows.Next() {
		var i Item
		var pr float64
		var gid int
		err := rows.Scan(&gid, &i.Name, &pr, &i.Quantity, &i.Category, &i.ItemId)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []GachaponItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		gci = append(gci, GachaponItem{Item: i,
			PullRate:   pr,
			GachaponID: gid,
		})
	}
	return gci, nil
}
