package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Println("Hello world")

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadRequest)
			return
		}

		log.Printf("Data: %s", d)
		fmt.Fprintf(rw, "Hello %s", d)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye world")
	})

	http.ListenAndServe("127.0.0.1:9090", nil) // ip:port, handler
}