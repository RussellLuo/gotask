package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RussellLuo/gotask"
)

type HTTPWorker struct {
	Registry gotask.Registry
	Addr     string
}

func (w *HTTPWorker) Work() error {
	http.HandleFunc("/", w.handle)
	log.Printf("Listening on %s", w.Addr)
	return http.ListenAndServe(w.Addr, nil)
}

func (w *HTTPWorker) handle(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sig := gotask.Signature{}
	if err := decoder.Decode(&sig); err != nil {
		log.Printf("err: %#v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := gotask.Process(w.Registry, &sig); err != nil {
		log.Printf("err: %#v", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
