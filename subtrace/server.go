package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

//go:embed index.html
var indexHtml []byte

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(indexHtml)
	if err != nil {
		http.Error(w, "Error writing data", http.StatusInternalServerError)
		return
	}
}

func streamDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for i := 0; i < 10000; i++ {
		t := time.Now().UnixNano() / int64(time.Millisecond) // unix time in ms
		y := rand.Intn(1024)
		row := make([]int, 512)
		for j := range row {
			row[j] = rand.Intn(2)
		}
		data := []interface{}{t, y, row}
		buf := bytes.NewBuffer(make([]byte, 0, 4096))
		err := json.NewEncoder(buf).Encode(data)
		if err != nil {
			http.Error(w, "Error encoding data", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(buf.Bytes())
		if err != nil {
			http.Error(w, "Error writing data", http.StatusInternalServerError)
			return
		}
		flusher.Flush()
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/zone/", streamDataHandler)
	fmt.Printf("======== Running on http://0.0.0.0:8080 ========\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
