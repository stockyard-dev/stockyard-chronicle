package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-chronicle/internal/server"
	"github.com/stockyard-dev/stockyard-chronicle/internal/store"
)

var version = "dev"

func main() {
	portFlag := flag.String("port", "", "HTTP port (overrides PORT env var)")
	dataFlag := flag.String("data", "", "Data directory (overrides DATA_DIR env var)")
	flag.Parse()

	port := *portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8760"
	}

	dataDir := *dataFlag
	if dataDir == "" {
		dataDir = os.Getenv("DATA_DIR")
	}
	if dataDir == "" {
		dataDir = "./chronicle-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("chronicle: %v", err)
	}
	defer db.Close()

	srv := server.New(db, server.DefaultLimits(), dataDir)

	fmt.Printf("\n  Chronicle v%s — Self-hosted event log\n", version)
	fmt.Printf("  Dashboard:  http://localhost:%s/ui\n", port)
	fmt.Printf("  API:        http://localhost:%s/api\n", port)
	fmt.Printf("  Data:       %s\n", dataDir)
	fmt.Printf("  Questions?  hello@stockyard.dev — I read every message\n\n")

	log.Printf("chronicle: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
