package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmandevv/blogging-platform-api/internal/blog"
	"github.com/dmandevv/blogging-platform-api/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func HandleGet(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("_id")
	if stringID == "" {
		http.Error(w, "Empty blog ID", http.StatusBadRequest)
		return
	}

	blogObjectID, err := bson.ObjectIDFromHex(stringID)
	if err != nil {
		http.Error(w, "Couldn't convert blog ID to ObjectID: "+err.Error(), http.StatusBadRequest)
		return
	}
	//retrieve from db
	if cfg.MongoClient != nil {

		collection := cfg.MongoClient.Database(cfg.MongoDB).Collection(cfg.MongoCollection)
		ctx, _ := context.WithTimeout(context.TODO(), cfg.MongoInsertTimeout)
		findResult := collection.FindOne(ctx, bson.M{"_id": blogObjectID})
		if findResult.Err() != nil {
			if findResult.Err() == mongo.ErrNoDocuments {
				http.Error(w, fmt.Sprintf("No document with this ID found: %s", blogObjectID), http.StatusNotFound)
				return
			}
			http.Error(w, "MongoDB had an issue looking up the blog: "+findResult.Err().Error(), http.StatusInternalServerError)
			return
		}

		var foundPost blog.BlogPost
		err := findResult.Decode(&foundPost)
		if err != nil {
			http.Error(w, "Failed to decode retrieved blog post: "+err.Error(), http.StatusInternalServerError)
			return
		}
		prettyPost, err := json.MarshalIndent(foundPost, "", " ")
		if err != nil {
			http.Error(w, "Blog found, but failed to marshal response", http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(prettyPost))
		return
	}
	http.Error(w, "Can't connect to MongoDB", http.StatusServiceUnavailable)
}

func HandleGetAll(cfg *config.Config, w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	var searchTerms []string
	for param := range queryParams {
		if param == "term" {
			searchTerms = append(searchTerms, queryParams.Get(param))
		}
	}

	var filter bson.D
	if len(searchTerms) == 0 {
		filter = bson.D{}
	} else {
		filter = bson.D{
			bson.E{
				Key: "$or",
				Value: bson.A{
					bson.M{"title": bson.M{"$in": searchTerms}},
					bson.M{"content": bson.M{"$in": searchTerms}},
					bson.M{"category": bson.M{"$in": searchTerms}},
				},
			},
		}
	}

	//retrieve all from db
	if cfg.MongoClient != nil {
		collection := cfg.MongoClient.Database(cfg.MongoDB).Collection(cfg.MongoCollection)
		ctx, _ := context.WithTimeout(context.TODO(), cfg.MongoInsertTimeout)

		findAllResult, err := collection.Find(ctx, filter)
		if err != nil {
			http.Error(w, "MongoDB had an issue looking up blogs: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//read results one by one
		defer findAllResult.Close(ctx)
		var allPosts []blog.BlogPost
		for findAllResult.Next(ctx) {
			var post blog.BlogPost
			if err := findAllResult.Decode(&post); err != nil {
				http.Error(w, "Error reading document: "+err.Error(), http.StatusInternalServerError)
				return
			}
			allPosts = append(allPosts, post)
		}

		if len(allPosts) == 0 {
			http.Error(w, "No blogs found", http.StatusOK)
			return
		}

		prettyPosts, err := json.MarshalIndent(allPosts, "", " ")
		if err != nil {
			http.Error(w, "Blogs found, but failed to marshal response", http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(prettyPosts))
		return
	}
	http.Error(w, "Can't connect to MongoDB", http.StatusServiceUnavailable)
}
