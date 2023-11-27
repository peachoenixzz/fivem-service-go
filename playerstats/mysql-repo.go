package playerstats

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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

type Player struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Money      int    `json:"money"`
	Bank       int    `json:"bank"`
	Total      int    `json:"total"`
	BlackMoney int    `json:"black_money"`
}

type PlayerItem struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Steel      string `json:"steel"`
	TotalSteel int    `json:"total_steel"`
}

type Vehicle struct {
	Owner       string      `json:"owner"`
	VehicleData string      `json:"vehicle_data"`
	Model       interface{} `json:"model"`
}

func (h Handler) VehicleByModel(ctx context.Context) error {
	logger := mlog.Logg
	logger.Info("prepare to make query Vehicle")
	stmtStr := `
    	SELECT 
			vehicle,
    		owner
		FROM 
			owned_vehicles
	`

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	count := 0
	var vehicles []Vehicle
	// Iterate over the rows
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.VehicleData, &vehicle.Owner); err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}

		// Parse the data JSON
		var values map[string]any
		if err := json.Unmarshal([]byte(vehicle.VehicleData), &values); err != nil {
			logger.Error("JSON UNMARSHAL ERR : ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		//fmt.Println(-204559)
		fmt.Println(values["model"])
		vehicle.Model = values["model"]
		// Update the totals
		if vehicle.Model == "-204559" {
			vehicles = append(vehicles, vehicle)
			count++
		}

	}

	if err := rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	fmt.Println("all vehicle model -204559 : ", count)

	for _, vehicle := range vehicles {
		fmt.Println(fmt.Sprintf("Player : %v %v ", vehicle.Owner, vehicle.Model))
	}

	if err = rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	return nil
}

func (h Handler) AllMoney(ctx context.Context) error {
	logger := mlog.Logg
	logger.Info("prepare to make query AllMoney")
	stmtStr := `
    	SELECT 
			accounts,
			firstname,
			lastname
		FROM 
			users
	`

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	totalBank := 0
	totalMoney := 0
	totalBlackMoney := 0
	count := 0
	var players []Player
	// Iterate over the rows
	for rows.Next() {
		var player Player
		var data string
		if err := rows.Scan(&data, &player.FirstName, &player.LastName); err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}

		// Parse the data JSON

		var values map[string]int
		if err := json.Unmarshal([]byte(data), &values); err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}

		// Update the totals
		totalBank += values["bank"]
		totalMoney += values["money"]
		totalBlackMoney += values["black_money"]
		player.Money = values["money"]
		player.Bank = values["bank"]
		player.Total = values["money"] + values["bank"]
		player.BlackMoney = values["black_money"]
		if (player.Total) >= 0 {
			players = append(players, player)
		}
		count++
	}

	if err := rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	// Print the totals
	//fmt.Println("Total bank:", totalBank)
	//fmt.Println("Total money:", totalMoney)
	//fmt.Println("Total green money before minus starter money :", (totalMoney+totalBank)-(count*3000))
	//fmt.Println("Total green money:", totalMoney+totalBank)
	//fmt.Println("Total black money:", totalBlackMoney)

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("PlayerData")
	if err != nil {
		logger.Error("Error creating the Excel sheet: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating the Excel sheet: ", err.Error())
	}

	// Write the header row
	row := sheet.AddRow()
	row.AddCell().SetValue("PlayerName")
	row.AddCell().SetValue("Bank")
	row.AddCell().SetValue("Money")
	row.AddCell().SetValue("Black Money")
	row.AddCell().SetValue("Total Money Player")
	row.AddCell().SetValue(fmt.Sprintf("Total green money before minus starter money : %v", (totalMoney+totalBank)-(count*3000)))
	row.AddCell().SetValue(fmt.Sprintf("Total Green Server : %v", totalMoney+totalBank))
	row.AddCell().SetValue(fmt.Sprintf("Total Black Server : %v", totalBlackMoney))

	// Write the data to the Excel sheet
	for _, player := range players {
		//fmt.Println(fmt.Sprintf("%d %s %s", player.Bank, player.FirstName, player.LastName))
		//fmt.Println(fmt.Sprintf("%d %s %s", player.Money, player.FirstName, player.LastName))
		//fmt.Println(fmt.Sprintf("%d %s %s", player.BlackMoney, player.FirstName, player.LastName))
		fmt.Println(fmt.Sprintf("%d %s %s", player.Total, player.FirstName, player.LastName))
		row := sheet.AddRow()
		row.AddCell().SetValue(player.FirstName + " " + player.LastName)
		row.AddCell().SetValue(fmt.Sprintf("%d", player.Bank))
		row.AddCell().SetValue(fmt.Sprintf("%d", player.Money))
		row.AddCell().SetValue(fmt.Sprintf("%d", player.BlackMoney))
		row.AddCell().SetValue(fmt.Sprintf("%d", player.Total))
	}

	// Save the Excel file
	err = file.Save(fmt.Sprintf("data-money-player-%v", time.Now()) + ".xlsx")
	if err != nil {
		logger.Error("Error saving the Excel file: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Error saving the Excel file: ", err.Error())
	}

	fmt.Println("Data imported successfully to player_data.xlsx")

	if err = rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	return nil
}

func (h Handler) ItemPlayer(ctx context.Context) error {
	logger := mlog.Logg
	logger.Info("prepare to make query AllMoney")
	stmtStr := `
    	SELECT 
			inventory,
			firstname,
			lastname
		FROM 
			users
	`

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	steelProcess := 0
	//count := 0
	//var pis []PlayerItem
	// Iterate over the rows
	for rows.Next() {
		var pi PlayerItem
		var data string
		if err := rows.Scan(&data, &pi.FirstName, &pi.LastName); err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}

		// Parse the data JSON

		if data != "[]" {
			var values map[string]int
			if err := json.Unmarshal([]byte(data), &values); err != nil {
				logger.Error("Database Error : ", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
			}

			// Update the totals
			steelProcess += values["steel_pro_1"]
		}

	}

	fmt.Println(steelProcess)

	if err := rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	return nil
}

func (h Handler) ItemVault(ctx context.Context) error {
	logger := mlog.Logg
	logger.Info("prepare to make query AllMoney")
	stmtStr := `
    	SELECT 
			items
		FROM 
			nc_vault_storage
	`

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	steelProcess := 0
	count := 0
	//var pis []PlayerItem
	// Iterate over the rows
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}

		if data != "[]" && data != "" {
			var values map[string]int
			if err := json.Unmarshal([]byte(data), &values); err != nil {
				logger.Error("Database Error : ", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
			}

			// Update the totals
			if values["steel_pro_1"] > 0 {
				steelProcess += values["steel_pro_1"]
				count++
			}
		}
	}
	fmt.Println(count)
	fmt.Println(steelProcess)

	if err := rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	return nil
}
