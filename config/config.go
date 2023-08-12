package config

import (
	"os"
	"strconv"
)

type env func(key string) string

type Cfg struct {
	getEnv env
}

func New() *Cfg {
	return &Cfg{getEnv: os.Getenv}
}

type Config struct {
	Server            Server
	FeatureFlag       FeatureFlag
	DBConnection      string
	MongoDBConnection string
	MySQLDBConnection string
}

type Server struct {
	Hostname string
	Port     int
}

type FeatureFlag struct {
	IsLimitMaxBalanceOnCreate bool `json:"isLimitMaxBalanceOnCreate"`
}

const (
	cHostname = "HOSTNAME"
	cPort     = "PORT"

	cFlagIsLimitMaxBalanceOnCreate = "FLAG_IS_LIMIT_MAX_SPEND_ON_CREATE"
	cDBConnection                  = "DB_CONNECTION"
)

const (
	dPort = 1323
	//dDBConnection     = "postgresql://postgres:password@localhost:5432/banking?sslmode=disable"
	dDBConnection     = "postgresql://postgres:1234@localhost:5432?sslmode=disable"
	mongoDBConnection = "mongodb://localhost:27017" //mongodb://fivemlogs:isylzjkoshkm001@mongodb/fivem-logs
	mySQLDBConnection = "peachoenixz:petuyio001@tcp(103.212.181.194:3306)/es_extended_feature"
	//mySQLDBConnection = "doraemonfivem:Doraemon001FiveM@tcp(103.212.181.194:3306)/es_extended"

	//mongoDBConnection = "mongodb://fivemlogs:isylzjkoshkm001@localhost/fivem-logs" //mongodb://fivemlogs:isylzjkoshkm001@mongodb/fivem-logs
	//mySQLDBConnection = "doraemonfivem:Doraemon001FiveM@tcp(localhost:3306)/"
)

func (c *Cfg) All() Config {
	return Config{
		Server: Server{
			Hostname: c.envString(cHostname, ""),
			Port:     c.envInt(cPort, dPort),
		},
		FeatureFlag: FeatureFlag{
			IsLimitMaxBalanceOnCreate: c.envBool(cFlagIsLimitMaxBalanceOnCreate, false),
		},
		DBConnection:      c.envString(cDBConnection, dDBConnection),
		MongoDBConnection: c.envString(cDBConnection, mongoDBConnection),
		MySQLDBConnection: c.envString(cDBConnection, mySQLDBConnection),
	}
}

func (c *Cfg) SetEnvGetter(overrideEnvGetter env) {
	c.getEnv = overrideEnvGetter
}

func (c *Cfg) envString(key, defaultValue string) string {
	val := c.getEnv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func (c *Cfg) envInt(key string, defaultValue int) int {
	v := c.getEnv(key)

	val, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}

	return val
}

func (c *Cfg) envBool(key string, defaultValue bool) bool {
	v := c.getEnv(key)

	val, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}

	return val
}
