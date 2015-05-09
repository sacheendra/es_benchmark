package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	es "github.com/mattbaird/elastigo/lib"
	"github.com/satori/go.uuid"
	"github.com/synful/grammar"
)

var esClient *es.Conn
var dummy *grammar.Grammar
var wg sync.WaitGroup

func main() {
	numThreads, _ := strconv.Atoi(os.Args[1])

	esClient = es.NewConn()
	esClient.SetPort("9200")
	esClient.SetHosts([]string{"localhost"})

	file, err := os.Open("./grammar/wellformedGrammar.txt")
	if err != nil {
		log.Fatalln("Error opening file", err)
	}

	dummy, err = grammar.New(file)
	if err != nil {
		log.Fatalln("Error generating speaker", err)
	}

	for i := 0; i < numThreads; i++ {
		go indexDocuments()
		wg.Add(1)
	}
	wg.Wait()
}

func indexDocuments() {
	for {
		indexDocument()
	}
}

func indexDocument() {
	id := uuid.NewV1().String()
	var timestamp int64 = time.Now().UnixNano()

	buf := new(bytes.Buffer)
	dummy.Speak(buf)

	body := map[string]interface{}{
		"upsert":          map[string]interface{}{},
		"scripted_upsert": true,
		"lang":            "groovy",
		"script_id":       "updateDocument",
		"params": map[string]interface{}{
			"timestamp": timestamp,
			"update": map[string]interface{}{
				"message": buf.String(),
			},
		},
	}

	_, err := esClient.Update("appbase", "bench1", id, nil, body)
	if err != nil {
		log.Println(err)
	}
}
