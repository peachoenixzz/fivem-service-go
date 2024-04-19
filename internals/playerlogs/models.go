package playerlogs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type RequestInsert struct {
	Event     string    `json:"event"`
	Content   string    `json:"content"`
	Source    int       `json:"source"`
	Color     string    `json:"color"`
	Options   Options   `json:"options"`
	Image     string    `json:"image"`
	Timestamp time.Time `json:"timestamp"`
	Player    Player    `json:"player"`
	Hardware  []string  `json:"hardware"`
}

type Options struct {
	Public    bool `json:"public"`
	Important bool `json:"important"`
}

type Player struct {
	Name        string      `json:"name"`
	Identifiers Identifiers `json:"identifiers"`
	Steam       PlayerSteam `json:"steam"`
}

type Identifiers struct {
	Ip       string `json:"ip"`
	Steam    string `json:"steam"`
	Discord  string `json:"discord"`
	License  string `json:"license"`
	License2 string `json:"license2"`
}

type PlayerSteam struct {
	Id     int    `json:"id"`
	Avatar string `json:"avatar"`
	Url    string `json:"url"`
}

type RequestCustomLog struct {
	DiscordID string `json:"discord_id,omitempty"`
	Event     string `json:"event,omitempty"`
	Begin     string `json:"begin"`
	Until     string `json:"until"`
	Regex     string `json:"regex,omitempty"`
}
