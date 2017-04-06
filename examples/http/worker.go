package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RussellLuo/gotask"
)

type HTTPWorker struct {
	Registry map[string]gotask.Constructor
	Addr     string
}

func (w *HTTPWorker) Start() error {
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
	err := decoder.Decode(&sig)
	if err != nil {
		log.Printf("err: %#v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = gotask.Process(w.Registry, &sig)
	if err != nil {
		log.Printf("err: %#v", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
