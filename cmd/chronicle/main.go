package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-chronicle/internal/server";"github.com/stockyard-dev/stockyard-chronicle/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8760"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./chronicle-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("chronicle: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Chronicle — Self-hosted event log\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,dataDir)
log.Printf("chronicle: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
