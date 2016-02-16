package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/sacheendra/es_benchmark/Godeps/_workspace/src/github.com/cenkalti/backoff"
	"github.com/sacheendra/es_benchmark/Godeps/_workspace/src/github.com/gosuri/uilive"
	"github.com/sacheendra/es_benchmark/Godeps/_workspace/src/github.com/joshlf/grammar"
	"github.com/sacheendra/es_benchmark/Godeps/_workspace/src/gopkg.in/olivere/elastic.v3"
)

var dummy *grammar.Grammar
var wg sync.WaitGroup
var numRequests, successfulRequests uint64

func main() {
	esURL := os.Getenv("ES_URL")
	sniffString := os.Getenv("ES_SNIFF")
	numClientsString := os.Getenv("NUM_CLIENTS")
	numThreadsString := os.Getenv("NUM_THREADS")

	if len(os.Args) >= 2 && os.Args[1] != "" {
		esURL = os.Args[1]
	}

	if len(os.Args) >= 3 && os.Args[2] != "" {
		numClientsString = os.Args[2]
	}

	if len(os.Args) >= 4 && os.Args[3] != "" {
		numThreadsString = os.Args[3]
	}

	if sniffString == "" {
		sniffString = "false"
	}

	sniff, err := strconv.ParseBool(sniffString)
	if err != nil {
		log.Fatalln(err)
	}

	numClients, err := strconv.Atoi(numClientsString)
	if err != nil {
		log.Fatalln(err)
	}

	numThreads, err := strconv.Atoi(numThreadsString)
	if err != nil {
		log.Fatalln(err)
	}

	numRequests = 0
	successfulRequests = 0

	go outputToTerminal()

	clients := make([]*elastic.Client, numClients)
	for i := 0; i < numClients; i++ {
		client, err := elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(sniff))
		if err != nil {
			log.Fatalln("Error initializing client: ", err)
		}

		clients[i] = client
	}

	for i := 0; i < numClients; i++ {
		for j := 0; j < numThreads; j++ {
			go indexDocuments(clients[i])
			wg.Add(1)
		}
	}

	wg.Wait()
}

func outputToTerminal() {
	writer := uilive.New()
	writer.Start()

	var requests, successfulReqs uint64
	for {
		requests = atomic.SwapUint64(&numRequests, 0)
		successfulReqs = atomic.SwapUint64(&successfulRequests, 0)

		fmt.Fprintf(writer, `
Requests per second: %d
Successful requests per second: %d
`, requests, successfulReqs)

		time.Sleep(1 * time.Second)
	}
}

func getSentenceGenerator() *grammar.Grammar {
	file, err := os.Open("./grammar/wellformedGrammar.txt")
	if err != nil {
		log.Fatalln("Error opening file", err)
	}

	dummy, err = grammar.New(file)
	if err != nil {
		log.Fatalln("Error generating speaker", err)
	}

	return dummy
}

func indexDocuments(client *elastic.Client) {
	var err error

	dummy := getSentenceGenerator()

	buf := new(bytes.Buffer)
	dummy.Speak(buf)

	b := NewExponentialBackOff()
	b.Reset()
	for {
		_, err = client.Index().Index("testindex").Type("bench1").BodyJson(map[string]interface{}{
			"message": buf.String(),
		}).Do()

		atomic.AddUint64(&numRequests, 1)

		if err == nil {
			atomic.AddUint64(&successfulRequests, 1)
		}

		if err != nil {
			elastic_err, ok := err.(*elastic.Error)
			if !ok {
				if !(strings.Contains(err.Error(), "http") || strings.Contains(err.Error(), "no Elasticsearch")) {
					log.Fatalln("Other error while making request: ", err)
				}
			}

			if elastic_err != nil && elastic_err.Status != 429 {
				log.Fatalln("Elasticsearch Client error: ", elastic_err.Error())
			}

			if skip := b.NextBackOff(); skip == Stop {
				b.Reset()
			} else {
				time.Sleep(skip)
			}
		}
	}
}
