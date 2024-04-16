package discordbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"log"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mysqlDB}
}

func (h Handler) handleEICommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	query := `
        SELECT i.label, DATE_FORMAT(ie.expire_timestamp, '%Y-%m-%d %H:%i:%s') 
        FROM users u
        INNER JOIN items_expire ie ON ie.player_id = u.identifier
        INNER JOIN items i ON i.name = ie.item_name
        WHERE ie.player_id = ?;
    `

	user, err := s.User(m.Author.ID)
	if err != nil || user.Bot {
		return // Ignore if it's a bot or user fetch error
	}

	args := []interface{}{
		user.ID,
	}

	channel, err := s.UserChannelCreate(user.ID)
	if err != nil {
		log.Printf("Error creating DM channel: %v", err)
		return
	}

	rows, err := h.MysqlDB.Query(query, args...) // Replace with dynamic ID if needed
	if err != nil {
		s.ChannelMessageSend(channel.ID, fmt.Sprintf("Error querying database: %v", err))
		return
	}
	defer rows.Close()

	var is []Item
	for rows.Next() {
		var i Item
		err := rows.Scan(&i.ItemName, &i.ExpireDate)
		if err != nil {
			s.ChannelMessageSend(channel.ID, fmt.Sprintf("Error reading row: %v", err))
			return
		}
		is = append(is, i)
	}
	for _, i := range is {
		response := fmt.Sprintf("ชื่อไอเทม : %s, วันหมดอายุ : %v\n", i.ItemName, i.ExpireDate)
		s.ChannelMessageSend(channel.ID, response)
	}
}

func (h Handler) handleEVCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	query := `
        SELECT ov.plate, DATE_FORMAT(ov.expire_date, '%Y-%m-%d %H:%i:%s') as expire_date
        FROM users u
        INNER JOIN owned_vehicles ov ON ov.owner = u.identifier
        WHERE ov.expire_date IS NOT NULL
        AND ov.owner = ?;
    `
	user, err := s.User(m.Author.ID)
	if err != nil || user.Bot {
		return // Ignore if it's a bot or user fetch error
	}

	args := []interface{}{
		user.ID,
	}

	channel, err := s.UserChannelCreate(user.ID)
	if err != nil {
		log.Printf("Error creating DM channel: %v", err)
		return
	}

	rows, err := h.MysqlDB.Query(query, args...) // Replace with dynamic ID if needed
	if err != nil {
		s.ChannelMessageSend(channel.ID, fmt.Sprintf("Error querying database: %v", err))
		return
	}
	defer rows.Close()

	var vehs []Vehicle
	for rows.Next() {
		var veh Vehicle
		err := rows.Scan(&veh.Plate, &veh.ExpireDate)
		if err != nil {
			s.ChannelMessageSend(channel.ID, fmt.Sprintf("Error reading row: %v", err))
			return
		}
		vehs = append(vehs, veh)
	}
	for _, v := range vehs {
		response := fmt.Sprintf("ทะเบียน : %s, วันหมดอายุ : %v\n", v.Plate, v.ExpireDate)
		s.ChannelMessageSend(channel.ID, response)
	}
}
