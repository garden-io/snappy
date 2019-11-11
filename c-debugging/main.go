package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	called := time.Unix(0, 1573314548127576044)
	now := time.Now()
	diff := now.Sub(called)
	fmt.Println(diff.Round(time.Millisecond))

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
