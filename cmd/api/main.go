package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	fmt.Println("listening on :4000")
	http.ListenAndServe(":4000", nil)
}
