package config

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Config struct {
	MongoClient        *mongo.Client
	MongoDB            string
	MongoCollection    string
	MongoInsertTimeout time.Duration
	Host               string
	Port               int
}
