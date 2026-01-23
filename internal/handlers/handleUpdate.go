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

func HandleUpdate(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("_id")
	if stringID == "" {
		http.Error(w, "Empty or invalid blog ID", http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON into a BlogPost struct
	var updatedPost blog.BlogPost
	err = json.Unmarshal(body, &updatedPost)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if updatedPost.Title == "" || updatedPost.Content == "" {
		http.Error(w, "Title and Content are required", http.StatusBadRequest)
		return
	}

	blogObjectID, err := bson.ObjectIDFromHex(stringID)
	if err != nil {
		http.Error(w, "Couldn't convert blog ID to ObjectID: "+err.Error(), http.StatusBadRequest)
		return
	}

	//check if db is running
	if cfg.MongoClient != nil {
		collection := cfg.MongoClient.Database(cfg.MongoDB).Collection(cfg.MongoCollection)
		ctx, _ := context.WithTimeout(context.TODO(), cfg.MongoInsertTimeout)
		findResult := collection.FindOne(ctx, bson.M{"_id": blogObjectID})
		if findResult.Err() != nil {
			http.Error(w, "Failed to find blog post in MongoDB: "+findResult.Err().Error(), http.StatusNotFound)
			return
		}

		var oldPost blog.BlogPost
		err := findResult.Decode(&oldPost)
		if err != nil {
			http.Error(w, "Failed to decode original blog: "+err.Error(), http.StatusInternalServerError)
			return
		}

		updatedPost.ID = oldPost.ID
		updatedPost.CreatedAt = oldPost.CreatedAt
		updatedPost.UpdatedAt = time.Now()

		_, err = collection.ReplaceOne(context.TODO(), bson.M{"_id": blogObjectID}, updatedPost)
		if err != nil {
			http.Error(w, "Failed to update blog post in MongoDB: "+err.Error(), http.StatusInternalServerError)
			return
		}

		prettyPost, err := json.MarshalIndent(updatedPost, "", " ")
		if err != nil {
			http.Error(w, "Blog created and saved to database, but failed to marshal response", http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(prettyPost))
		return
	}

	http.Error(w, "Can't connect to MongoDB", http.StatusServiceUnavailable)
}
