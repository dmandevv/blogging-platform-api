package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dmandevv/blogging-platform-api/internal/config"
	"github.com/dmandevv/blogging-platform-api/internal/handlers"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {

	cfg := &config.Config{
		MongoClient:        connectMongoDB(),
		MongoDB:            "blog-cluster",
		MongoCollection:    "posts",
		MongoInsertTimeout: time.Second * 5,
		Host:               "localhost",
		Port:               8080,
	}
	defer func() {
		err := cfg.MongoClient.Disconnect(context.TODO())
		if err != nil {
			panic(err)
		}
	}()

	mux := http.NewServeMux()

	//get
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAll(cfg, w, r)
	})
	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAll(cfg, w, r)
	})
	mux.HandleFunc("GET /posts/{_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGet(cfg, w, r)
	})

	//create
	mux.HandleFunc("POST /posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCreate(cfg, w, r)
	})

	//update
	mux.HandleFunc("PUT /posts/{_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUpdate(cfg, w, r)
	})

	//delete
	mux.HandleFunc("DELETE /posts/{_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleDelete(cfg, w, r)
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), mux)

}

func connectMongoDB() *mongo.Client {

	url, err := url.Parse("mongodb+srv://blogAdmin:zu@gbPDUYBXi3Kr@blog-cluster.ycanw9t.mongodb.net/?appName=blog-cluster")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connecting to MongoDB at:", url.String())
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url.String()).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	//ping server
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}
