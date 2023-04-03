package playerlogin

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kkgo-software-engineering/workshop/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type Handler struct {
	Cfg        config.FeatureFlag
	PostgresDB *sql.DB
	MongoDB    *mongo.Client
}

func New(cfgFlag config.FeatureFlag, postgresDB *sql.DB, mongoDB *mongo.Client) *Handler {
	return &Handler{cfgFlag, postgresDB, mongoDB}
}

func (h Handler) InsertMLog(req Request) (Message, error) {
	col := h.MongoDB.Database("fivem-logs").Collection("policelogs")
	_, err := col.InsertOne(context.Background(), req)
	if err != nil {
		return Message{Status: http.StatusInternalServerError, Message: "Database Failed"}, err
	}
	return Message{Status: http.StatusCreated, Message: "Created Success"}, nil
}

func (h Handler) FiveMLog() ([]Response, error) {
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100)
	col := h.MongoDB.Database("fivem-logs").Collection("policelogs")
	cur, err := col.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	var res []Response
	if err := cur.All(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (h Handler) LogCaseEventAndSteamID(steamid string, event string) ([]Response, error) {
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100)
	f := bson.M{
		"player.identifiers.steam": fmt.Sprintf("steam:%s", steamid),
		"event":                    event,
	}
	col := h.MongoDB.Database("fivem-logs").Collection("policelogs")
	cur, err := col.Find(context.Background(), f, opts)
	if err != nil {
		return nil, err
	}

	var res []Response
	if err := cur.All(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (h Handler) LogAllEventAndSteamID(steamid string) ([]Response, error) {
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100)
	fmt.Println(fmt.Sprintf("steam:%s", steamid))
	f := bson.M{
		"player.identifiers.steam": fmt.Sprintf("steam:%s", steamid),
	}
	col := h.MongoDB.Database("fivem-logs").Collection("policelogs")
	cur, err := col.Find(context.Background(), f, opts)
	if err != nil {
		return nil, err
	}

	var res []Response
	if err := cur.All(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (h Handler) LogCaseEventAll(event string) ([]Response, error) {
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100)
	col := h.MongoDB.Database("fivem-logs").Collection("policelogs")

	cur, err := col.Find(context.Background(), bson.M{"event": event}, opts)
	if err != nil {
		return nil, err
	}

	var res []Response
	if err := cur.All(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}
