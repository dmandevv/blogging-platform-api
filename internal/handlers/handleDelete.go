package handlers

import (
	"context"
	"net/http"

	"github.com/dmandevv/blogging-platform-api/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func HandleDelete(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("_id")
	if stringID == "" {
		http.Error(w, "Empty or invalid blog ID", http.StatusBadRequest)
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
		deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": blogObjectID})
		if err != nil {
			http.Error(w, "Failed to delete blog post in MongoDB: "+err.Error(), http.StatusNotFound)
			return
		}

		if deleteResult.DeletedCount == 0 {
			http.Error(w, "Failed to find blog", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Blog deleted successfully"))
		return
	}

	http.Error(w, "Can't connect to MongoDB", http.StatusServiceUnavailable)
}
