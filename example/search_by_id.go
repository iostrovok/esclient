package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/iostrovok/esclient"
	"github.com/olivere/elastic/v7"
)

var url, index, typeDoc, reqVal string

var countGoroutine int = 10
var printLock sync.RWMutex

func init() {

	printLock = sync.RWMutex{}

	log.SetFlags(0)

	flag.StringVar(&url, "url", "", "Elasticsearch URL")
	flag.StringVar(&index, "index", "", "Elasticsearch index")
	flag.StringVar(&typeDoc, "type", "", "Elasticsearch type")
	flag.StringVar(&reqVal, "req", "", "Searching data or id")

	flag.Parse()
}

func main() {

	// Create an Elasticsearch client
	client, err := esclient.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < countGoroutine; i++ {
		wg.Add(1)
		go runID(i, wg, client)
	}

	wg.Wait()
}

func runID(i int, wg *sync.WaitGroup, client esclient.IConn) {

	defer wg.Done()

	cl := client.Open(true)

	result, err := cl.Get().Get().
		Index(index).
		Type(typeDoc).
		Id(reqVal).
		Do(context.Background())

	// We don't want to mash several outputs for readability.
	printLock.Lock()
	defer printLock.Unlock()

	log.Printf("\n----------- %d -----------\n", i)
	log.Printf("Debug.Request ===> \n%s\n", cl.Debug().Request())
	log.Printf("Debug.Response ===> \n%s\n", cl.Debug().Response())
	log.Printf("----------------------\n")

	if err != nil {
		log.Println(err)
		return
	}

	if result.Found {
		log.Printf("Got document %s in version %d from index %s, type %s\n", result.Id, result.Version, result.Index, result.Type)
	}

	log.Println("Finished succeeded")
}
