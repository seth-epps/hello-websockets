package main

import (
	"context"
	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handleGet)
	r.Handle("/ws", http.HandlerFunc(handleWS))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	res := Response{"Hello!"}
	json.NewEncoder(w).Encode(res)
}

func handleWS(w http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("error accepting websocket connection: %v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")
	for {
		ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Printf("error reading from websocket connection: %v", err)
			return
		}
		log.Printf("received: %v\n", v)
		wsjson.Write(ctx, conn, Response{fmt.Sprintf("received: %v", v)})
	}
}
