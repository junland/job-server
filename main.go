package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Config does stuff
type Config struct {
	Workers         int
	ServeAddress    string
	WorkSpaceDir    string
	StartingHistory int
	CurrentHistory  int
}

var (
	//NWorkers describes number of workers to start
	NWorkers = flag.Int("n", 4, "The number of workers to start")

	//HTTPAddr describes what address to listen on.
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")

	//WorkspaceDir what directory to write and read the build history.
	WorkspaceDir = flag.String("workspace", "./workspace", "directory for build history.")
)

func main() {
	// Parse the command-line flags.
	flag.Parse()

	c := Config{Workers: *NWorkers, ServeAddress: *HTTPAddr, WorkSpaceDir: *WorkspaceDir, StartingHistory: 0, CurrentHistory: 0}

	c.StartingHistory = FindHistory(c.WorkSpaceDir)

	c.CurrentHistory = c.StartingHistory

	if _, err := os.Stat(*WorkspaceDir); os.IsNotExist(err) {
		fmt.Println("Workspace directory does not exist. Creating...")
		os.Mkdir(*WorkspaceDir, 0777)
	}

	fmt.Println("Starting the dispatcher")

	StartDispatcher(*NWorkers, *WorkspaceDir)

	middleware := alice.New(middlewareOne)

	router := httprouter.New()

	router.Handler("POST", "/work", middleware.ThenFunc(c.Collector))

	router.Handler("POST", "/stop", middleware.ThenFunc(c.StopWorker))

	fmt.Println("HTTP server listening on ", *HTTPAddr)

	log.Fatal(http.ListenAndServe(*HTTPAddr, router))
}

// middlewareOne is just a template middleware.
func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}
