package blog

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BlogPost struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string        `bson:"title" json:"title"`
	Content   string        `bson:"content" json:"content"`
	Category  string        `bson:"category" json:"category"`
	Tags      []string      `bson:"tags" json:"tags"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
