package mlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

const key = "logger"

var Logg *zap.Logger

type ElasticWriter struct {
	mu sync.Mutex
}

type LogData struct {
	Level   string
	Time    time.Time
	Message string
}

type ElasticsearchWriter struct {
	Client *elasticsearch.Client
	Name   string
}

func NewElasticsearchWriter(client *elasticsearch.Client, name string) *ElasticsearchWriter {
	return &ElasticsearchWriter{
		Client: client,
		Name:   name,
	}
}

func (ew *ElasticsearchWriter) Write(p []byte) (n int, err error) {
	var logEntry map[string]interface{}
	if err := json.Unmarshal(p, &logEntry); err != nil {
		return 0, err
	}

	if t, ok := logEntry["timestamp"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, t); err == nil {
			logEntry["time"] = parsedTime.Format(time.RFC3339) // Properly format or adjust the time
		} else {
			logEntry["time"] = time.Now().Format(time.RFC3339) // Use current time if parsing fails
		}
		delete(logEntry, "timestamp") // Remove the original timestamp field
	} else {
		logEntry["time"] = time.Now().Format(time.RFC3339) // Use current time if not present
	}

	// Convert modified log data back to JSON
	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		return 0, err
	}

	//fmt.Println(strings.NewReader(string(logJSON)))
	location, err := time.LoadLocation("Asia/Jakarta") // Specify the desired time zone
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	currentDate := time.Now().In(location).Format("2006-01-02")
	fmt.Println(currentDate) // '2018-01-02
	fmt.Println(ew.Name)     // 'fivem
	indexName := fmt.Sprintf("%s-logs-%s", ew.Name, currentDate)
	req := esapi.IndexRequest{
		Index:   indexName,
		Body:    bytes.NewReader(logJSON),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), ew.Client)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return len(p), nil
}

func L(c echo.Context) *zap.Logger {
	switch logger := c.Get(key).(type) {
	case *zap.Logger:
		return logger
	default:
		return zap.NewNop()
	}
}

func SetupLogger(es *elasticsearch.Client, name string) (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zap.InfoLevel

	writer := NewElasticsearchWriter(es, name)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writer)),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= level }),
	)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger) // Optional: Replace global logger

	Logg = logger
	return logger, nil
}
