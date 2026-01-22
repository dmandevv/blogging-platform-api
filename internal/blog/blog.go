package blog

import "time"

type BlogPost struct {
	ID        int       `bson:"_id,omitempty"`
	Title     string    `bson:"title" json:"title"`
	Content   string    `bson:"content" json:"content"`
	Category  string    `bson:"category" json:"category"`
	Tags      []string  `bson:"tags" json:"tags"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
