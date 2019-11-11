package main

import (
	"fmt"
	_ "github.com/dgraph-io/badger"
	_ "github.com/influxdata/influxdb"
	"net/http"
	"time"
)

func main() {
	called := time.Unix(0, 1573314099432078414)
	now := time.Now()
	diff := now.Sub(called)
	fmt.Println(diff.Round(time.Millisecond))

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
