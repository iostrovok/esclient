package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/iostrovok/esclient"
)

var (
	url, index, id string
	countGoroutine = 10
	printLock      sync.RWMutex
)

func init() {

	printLock = sync.RWMutex{}

	flag.StringVar(&url, "url", "", "Elasticsearch URL")
	flag.StringVar(&index, "index", "", "Elasticsearch index")
	flag.StringVar(&id, "id", "", "Searching id")

	flag.Parse()
}

// InitClient prepares instance of elastic client
func InitClient(url []string, user, password string) {

}

func main() {
	log.Printf("Start with:\n")
	log.Printf("url: %s\n", url)
	log.Printf("index: %s\n", index)
	log.Printf("id: %s\n\n", id)

	options := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetSniff(false), // Disabled, see https://github.com/olivere/elastic/issues/312
		elastic.SetHealthcheck(false),
	}

	connection, err := esclient.NewClient(options...)
	if err != nil {
		log.Println("connection, err := esclient.NewClient(options...)")
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	connection.SetLogger(log.New(os.Stderr, "INFO: ", log.Lshortfile))
	connection.SniffTimeout(1 * time.Second)
	connection.Sniff(ctx)

	time.Sleep(20 * time.Second)

	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)
		runID(i, connection)
	}

	cancel()
}

func runID(i int, client esclient.IConn) {
	cl := client.Open(true)

	result, err := cl.Get().Get().
		Index(index).
		Id(id).
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
