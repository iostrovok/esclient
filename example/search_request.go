package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"sync"

	elastic "github.com/olivere/elastic/v7"

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
	client, err := esclient.Dial(elastic.SetURL(url), elastic.SetSniff(false))
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
	log.Printf("Debug.Request ===> \n%s\n", cl.Debug().Request())
	log.Printf("Debug.Response ===> \n%s\n", cl.Debug().Response())
	log.Printf("----------------------\n")
	
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
