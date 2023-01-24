package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/olivere/elastic/v7"

	"github.com/iostrovok/esclient"
)

var (
	url, index, reqVal string
	countGoroutine     = 10
	printLock          sync.RWMutex
)

func init() {

	printLock = sync.RWMutex{}

	log.SetFlags(0)

	flag.StringVar(&url, "url", "", "Elasticsearch URL")
	flag.StringVar(&index, "index", "", "Elasticsearch index")
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

	cl := client.Open(true, context.Background())

	result, err := cl.Get().Get().
		Index(index).
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
