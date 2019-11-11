package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/dgraph-io/badger"
	_ "github.com/influxdata/influxdb"
)

func main() {
	called := time.Unix(0, 1573315163434929134)
	now := time.Now()
	diff := now.Sub(called)
	fmt.Println(diff.Round(time.Millisecond))

	// _ = packr.NewBox("./import")

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
