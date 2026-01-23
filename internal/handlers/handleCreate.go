package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/dmandevv/blogging-platform-api/internal/blog"
	"github.com/dmandevv/blogging-platform-api/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func HandleCreate(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON into a BlogPost struct
	var newPost blog.BlogPost
	err = json.Unmarshal(body, &newPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newPost.Title == "" || newPost.Content == "" {
		http.Error(w, "Title and Content are required", http.StatusBadRequest)
		return
	}

	//insert into db
	if cfg.MongoClient != nil {
		newPost.CreatedAt = time.Now()
		newPost.UpdatedAt = time.Now()

		collection := cfg.MongoClient.Database(cfg.MongoDB).Collection(cfg.MongoCollection)
		ctx, _ := context.WithTimeout(context.TODO(), cfg.MongoInsertTimeout)
		insertResult, err := collection.InsertOne(ctx, newPost)
		if err != nil {
			http.Error(w, "Failed to insert blog post into MongoDB: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if oid, ok := insertResult.InsertedID.(bson.ObjectID); ok {
			newPost.ID = oid
		} else {
			http.Error(w, "Failed to convert inserted ID to ObjectID", http.StatusInternalServerError)
			return
		}

		prettyPost, err := json.MarshalIndent(newPost, "", " ")
		if err != nil {
			http.Error(w, "Blog created and saved to database, but failed to marshal response", http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(prettyPost))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blog created BUT not saved to database: " + newPost.Title))
}
