package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/olivere/elastic"

	"github.com/iostrovok/esclient"
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

func runID(i int, wg *sync.WaitGroup, client esclient.IESClient) {

	defer wg.Done()

	cl := client.Open(esclient.ErrorAndDebug)
	defer cl.Close()

	result, err := cl.Get().Get().
		Index(index).
		Type(typeDoc).
		Id(reqVal).
		Do(context.Background())

	// We don't want to mash several outputs for readability.
	printLock.Lock()
	defer printLock.Unlock()

	log.Printf("\n----------- %d -----------\n", i)
	log.Printf("id: Error: %s\n", cl.GetError().Error())
	log.Printf("id: Status: %d\n", cl.GetError().Status())
	log.Printf("id: Code: %d\n", cl.GetError().Code())
	log.Printf("id: Reason: %s\n", cl.GetError().Reason())
	log.Printf("id: Type: %s\n", cl.GetError().Type())
	log.Printf("----------------------\n")

	if cl.GetError().Code() == esclient.Internal || cl.GetError().Code() == esclient.Unknown {
		log.Println(cl.GetError().Error(), "; ", cl.GetError().Reason(), ";", cl.GetError().Type())
	}

	if err != nil {
		log.Println(err)
		return
	}

	if result.Found {
		log.Printf("Got document %s in version %d from index %s, type %s\n", result.Id, result.Version, result.Index, result.Type)
	}

	log.Println("Finished succeeded")
}
