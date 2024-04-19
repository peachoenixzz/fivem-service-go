package elshandler

import "time"

type ElasticModel struct {
	DiscordID string      `json:"discord_id"`
	Date      time.Time   `json:"date"`
	Data      interface{} `json:"data"`
}
