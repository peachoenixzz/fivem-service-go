package fivemlogs

import (
	"context"
	"database/sql"
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
	col := h.MongoDB.Database("fivem-logs").Collection("fivemlogs")
	_, err := col.InsertOne(context.Background(), req)
	if err != nil {
		return Message{Status: http.StatusInternalServerError, Message: "Database Failed"}, err
	}
	return Message{Status: http.StatusCreated, Message: "Created Success"}, nil
}

func (h Handler) FiveMLog() ([]Response, error) {
	var res []Response
	findOptions := options.Find()
	findOptions.SetLimit(100)
	col := h.MongoDB.Database("fivem-logs").Collection("fivemlogs")

	cur, err := col.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return []Response{}, err
	}

	if err := cur.All(context.Background(), &res); err != nil {
		return []Response{}, err
	}
	if err != nil {
		return []Response{}, err
	}
	return res, nil
}

func (h Handler) LogByEventAndSteamID(steamid int, event string) ([]Response, error) {
	var res []Response
	findOptions := options.Find()
	findOptions.SetLimit(100)
	col := h.MongoDB.Database("fivem-logs").Collection("fivemlogs")

	cur, err := col.Find(context.Background(), bson.M{"$and": []bson.M{{"player.steam.id": steamid}, {"event": event}}}, findOptions)
	if err != nil {
		return []Response{}, err
	}

	if err := cur.All(context.Background(), &res); err != nil {
		return []Response{}, err
	}
	if err != nil {
		return []Response{}, err
	}
	return res, nil
}

func (h Handler) LogAllEventAndSteamID(steamid int) ([]Response, error) {
	var res []Response
	findOptions := options.Find()
	findOptions.SetLimit(100)
	col := h.MongoDB.Database("fivem-logs").Collection("fivemlogs")

	cur, err := col.Find(context.Background(), bson.M{"player.steam.id": steamid}, findOptions)
	if err != nil {
		return []Response{}, err
	}

	if err := cur.All(context.Background(), &res); err != nil {
		return []Response{}, err
	}
	if err != nil {
		return []Response{}, err
	}
	return res, nil
}

func (h Handler) LogCaseEventAll(event string) ([]Response, error) {
	var res []Response
	findOptions := options.Find()
	findOptions.SetLimit(100)
	col := h.MongoDB.Database("fivem-logs").Collection("fivemlogs")

	cur, err := col.Find(context.Background(), bson.M{"event": event}, findOptions)
	if err != nil {
		return []Response{}, err
	}

	if err := cur.All(context.Background(), &res); err != nil {
		return []Response{}, err
	}
	if err != nil {
		return []Response{}, err
	}
	return res, nil
}
