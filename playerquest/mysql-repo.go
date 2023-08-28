package playerquest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}

func (h Handler) QueryRequireQuest(c echo.Context, discordID string) (ResponseRequireQuestPlayer, error) {
	logger := mlog.Logg
	query := `
		SELECT count(tpw.discord_id),tlw.quantity , tpw.weight_level
              FROM TB_PLAYER_WEIGHT tpw
              INNER JOIN TB_LEVEL_WEIGHT tlw
 			  WHERE tlw.weight_level = tpw.weight_level +1 
              AND discord_id = ?
`

	// Create a prepared statement
	logger.Info("mysql prepare query Discord ID")
	stmt, err := h.MysqlDB.Prepare(query)
	if err != nil {
		logger.Error("sql error", zap.Error(err))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	logger.Info("query Row Discord Id")
	var rowCount int
	var res ResponseRequireQuestPlayer
	err = stmt.QueryRow(discordID).Scan(&rowCount, &res.Quantity, &res.WeightLevel)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return ResponseRequireQuestPlayer{}, err
	}
	logger.Info("after query row and ready to return Data")
	if rowCount == 0 {
		return ResponseRequireQuestPlayer{}, c.JSON(http.StatusNotFound, Message{Message: "not found discord id"})
	}
	return res, nil
}

func (h Handler) QueryQuestItem(ctx context.Context) ([]ResponseQuestItem, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	var items []ResponseQuestItem
	stmtStr := `SELECT name,rare FROM TB_ITEM_QUEST`

	//args := []interface{}{
	//	discordID,
	//	discordID,
	//	discordID,
	//}

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	for rows.Next() {
		var item ResponseQuestItem
		err := rows.Scan(&item.Name, &item.Rare)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []ResponseQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
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

func (h Handler) InsertSelectQuestItem(rsi []ResponseSelectedItem, discordID string) {
	logger := mlog.Logg
	logger.Info("prepare to make query Insert Quest Player")
	for _, item := range rsi {
		args := []interface{}{
			discordID,
			item.Name,
			item.Quantity,
		}
		stmtStr := "INSERT INTO TB_PLAYER_QUEST (discord_id, name, quantity) VALUES (?, ?, ? );"
		_, err := h.MysqlDB.Exec(stmtStr, args...)
		if err != nil {
			logger.Error("failed to insert quest item into database:", zap.Error(err))
		}
	}
	logger.Info("success Insert Quest Player")
}

func (h Handler) GetPlayerQuestItem(ctx context.Context, discordID string) ([]ResponsePlayerQuestItem, error) {
	logger := mlog.Logg
	var items []ResponsePlayerQuestItem
	stmtStr := `
			SELECT 
				subquery.name,
				quantity,
				i.label
			FROM (
				SELECT 
					discord_id,
					CASE
						WHEN (HOUR(created_date) BETWEEN 0 AND 5) THEN '0.00 - 6.00'
						WHEN (HOUR(created_date) BETWEEN 6 AND 11) THEN '6.00 - 12.00'
						WHEN (HOUR(created_date) BETWEEN 12 AND 17) THEN '12.00 - 18.00'
						WHEN (HOUR(created_date) BETWEEN 18 AND 23) THEN '18.00 - 0.00'
					END AS time_range,
					name,
					quantity,
					CASE
						WHEN COUNT(*) OVER (PARTITION BY discord_id, time_range) > 1 THEN 'already'
						ELSE 'none'
					END AS status
				FROM TB_PLAYER_QUEST
				WHERE discord_id = ?
				AND DATE(created_date) = CURDATE()
				AND status = 'in_progress'
			) AS subquery
			INNER JOIN items i 
			ON	i.name = subquery.name
			WHERE subquery.time_range = CASE
										  WHEN HOUR(NOW()) BETWEEN 0 AND 5 THEN '0.00 - 6.00'
										  WHEN HOUR(NOW()) BETWEEN 6 AND 11 THEN '6.00 - 12.00'
										  WHEN HOUR(NOW()) BETWEEN 12 AND 17 THEN '12.00 - 18.00'
										  WHEN HOUR(NOW()) BETWEEN 18 AND 23 THEN '18.00 - 0.00'
									   END
			ORDER BY discord_id, time_range, name;`

	logger.Info("mysql prepare query Discord ID")
	args := []interface{}{
		discordID,
	}

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponsePlayerQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponsePlayerQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	for rows.Next() {
		var item ResponsePlayerQuestItem
		err := rows.Scan(&item.ItemName, &item.Quantity, &item.LabelName)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []ResponsePlayerQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponsePlayerQuestItem{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	return items, nil
}

func (h Handler) GetStateQuest(ctx context.Context, discordID string) bool {
	logger := mlog.Logg
	logger.Info("prepare to make query check already get Quest Player")
	stmtStr := `SELECT
				count(discord_id) as rowCount
			FROM (
				SELECT
			discord_id,
			CASE
				WHEN HOUR(created_date) >= 0 AND HOUR(created_date) < 6 THEN "0.00 - 6.00"
				WHEN HOUR(created_date) >= 6 AND HOUR(created_date) < 12 THEN "6.00 - 12.00"
				WHEN HOUR(created_date) >= 12 AND HOUR(created_date) < 18 THEN "12.00 - 18.00"
				ELSE '18.00 - 0.00'
			END AS time_range,
			CASE
				WHEN COUNT(id) > 1 THEN 'already'
				ELSE 'none'
			END AS status
			FROM TB_PLAYER_QUEST
			WHERE discord_id = ?
			AND DATE(created_date) = CURDATE()
			AND status = 'in_progress'
			GROUP BY discord_id, time_range
			) AS subquery
			WHERE subquery.time_range = 
			CASE
				WHEN HOUR(NOW()) BETWEEN 0 AND 5 THEN "0.00 - 6.00"
				WHEN HOUR(NOW()) BETWEEN 6 AND 11 THEN "6.00 - 12.00"
				WHEN HOUR(NOW()) BETWEEN 12 AND 17 THEN "12.00 - 18.00"
				WHEN HOUR(NOW()) BETWEEN 18 AND 23 THEN "18.00 - 0.00"
			END
			GROUP BY discord_id, time_range
			ORDER BY discord_id, time_range`

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("sql error", zap.Error(err))
		return false
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	logger.Info("query player quest")
	var rowCount int
	err = stmt.QueryRow(discordID).Scan(&rowCount)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Info(fmt.Sprintf("player none quest in this time : %v", time.Now()))
		return true
	}
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return false
	}
	logger.Info("after query row and ready to return Data")
	logger.Info(fmt.Sprintf("player have quest in this time : %v", time.Now()))
	return false
}

func (h Handler) ResetQuest(tx *sql.Tx, discordID string) error {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	stmtStr := `

UPDATE TB_PLAYER_QUEST AS tpq
SET tpq.status = 'cancel'
WHERE EXISTS (
    SELECT 1
    FROM TB_PLAYER_QUEST AS subquery
    WHERE subquery.discord_id = ?
    AND DATE(subquery.created_date) = CURDATE()
    AND subquery.status = 'in_progress'
    AND (
        (HOUR(subquery.created_date) BETWEEN 0 AND 5 AND HOUR(NOW()) BETWEEN 0 AND 5)
        OR (HOUR(subquery.created_date) BETWEEN 6 AND 11 AND HOUR(NOW()) BETWEEN 6 AND 11)
        OR (HOUR(subquery.created_date) BETWEEN 12 AND 17 AND HOUR(NOW()) BETWEEN 12 AND 17)
        OR (HOUR(subquery.created_date) BETWEEN 18 AND 23 AND HOUR(NOW()) BETWEEN 18 AND 23)
    )
)
AND tpq.discord_id = ?
AND DATE(tpq.created_date) = CURDATE()
AND tpq.status = 'in_progress';
`

	r, err := tx.Exec(stmtStr, discordID, discordID)
	if err != nil {
		// If there is an error, rollback the transaction
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to retrieve affected rows: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("Row Affected Update cash point table", zap.Int64("row affected", rowsAffected))
	return nil
}
