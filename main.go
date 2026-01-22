package main

import (
	"context"
	"net/http"

	"github.com/dmandevv/blogging-platform-api/internal/config"
	"github.com/dmandevv/blogging-platform-api/internal/handlers"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	mongoClient, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:8080"))
	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	cfg := &config.Config{
		MongoClient: mongoClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("BLOG"))
	})
	mux.HandleFunc("POST /posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCreate(cfg, w, r)
	})

	http.ListenAndServe(":8080", mux)
}
