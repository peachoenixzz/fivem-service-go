package elshandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
	"time"
)

func InsertElasticsearch(els *elasticsearch.TypedClient, data interface{}, discordID string, index string, id string) {
	elsModel := ElasticModel{
		DiscordID: discordID,
		Date:      time.Now(),
		Data:      data,
	}
	jsonStr, err := json.Marshal(elsModel)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}
	fmt.Println(string(jsonStr))

	req := esapi.IndexRequest{
		Index:      "cashshop-logs",
		DocumentID: "discordID",
		Body:       strings.NewReader(string(jsonStr)),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), els)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatalf("Error indexing document: %s", res.String())
	}
}
