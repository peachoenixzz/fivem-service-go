package playerquest

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}
