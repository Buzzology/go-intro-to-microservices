package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}


func (h *Goodbye) ServeHTTP(rw http.ResponseWriter, req *http.Request){
	h.l.Println("Goodbye...")

	d, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Goodbye %s", d)
}