package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"sync"

	"github.com/olivere/elastic"

	"github.com/iostrovok/esclient"
)

var url, index, typeDoc, reqVal string
var sortField, reqField string

var countGoroutine int = 10
var printLock sync.RWMutex

func init() {

	printLock = sync.RWMutex{}

	log.SetFlags(0)

	flag.StringVar(&url, "url", "", "Elasticsearch URL")
	flag.StringVar(&index, "index", "", "Elasticsearch index")
	flag.StringVar(&typeDoc, "type", "", "Elasticsearch type")
	flag.StringVar(&reqVal, "req", "", "Searching data or id")
	flag.StringVar(&sortField, "sort", "", "Sorting field")
	flag.StringVar(&reqField, "field", "", "Searching field")

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

	q := elastic.NewMatchQuery(reqField, reqVal)
	sortBy := elastic.SortInfo{
		Field:     sortField,
		Ascending: true,
	}
	searchResult, err := cl.Get().Search().
		Index(index).
		Type(typeDoc).
		Query(q).
		SortWithInfo(sortBy).
		Size(1).
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

	log.Println(">>>>>>>>>>>>>>>>>>>\n", string(cl.GetDebug().Request()), "<<<<<<<<<<<<<<<<<<<<")

	if err != nil {
		log.Println(err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		one := map[string]interface{}{}
		err := json.Unmarshal(*hit.Source, &one)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(one)
	}

	log.Println("Finished succeeded")
}
