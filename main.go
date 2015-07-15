package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	. "github.com/cenkalti/backoff"
	"github.com/olivere/elastic"
	"github.com/synful/grammar"
)

var client *elastic.Client
var dummy *grammar.Grammar
var wg sync.WaitGroup

func main() {
	numThreads, _ := strconv.Atoi(os.Args[1])

	var err error
	client, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatalln(err)
	}

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
	b := NewExponentialBackOff()
	b.Reset()
	for {
		err := indexDocument()
		if err != nil {
			if skip := b.NextBackOff(); skip == Stop {
				log.Println(err)
				b.Reset()
			} else {
				time.Sleep(skip)
			}
		}
	}
}

func indexDocument() error {
	buf := new(bytes.Buffer)
	dummy.Speak(buf)

	_, err := client.Index().Index("appbase").Type("bench1").BodyJson(map[string]interface{}{
		"message": buf.String(),
	}).Do()

	return err
}
