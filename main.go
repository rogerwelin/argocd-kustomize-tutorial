package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var Version = "dev"

type Message struct {
	Msg     string `json:"msg"`
	Version string `json:"version"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "pong")
}

func helloArgo(w http.ResponseWriter, r *http.Request) {
	msg := &Message{Msg: "Hello ArgoCD!", Version: Version}
	out, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(out)
}

func main() {
	wait := time.Second * 3

	http.HandleFunc("/", helloArgo)
	http.HandleFunc("/ping", pingHandler)
	srv := &http.Server{
		Addr: ":8080",
	}

	// Non-block so we can catch SIGINT
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Application started successfully")

	// Graceful shutdown on SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive a OS signal
	<-c
	log.Println("Shutdown signal received...")
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down server...")
}
