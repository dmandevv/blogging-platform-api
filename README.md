# Blogging Platform API

A RESTful API for a personal blogging platform built with Go and MongoDB. This API provides full CRUD functionality for managing blog posts with support for categorization and tagging.

This project is an exercise from the roadmap.sh project: [blogging-platform-api](https://roadmap.sh/projects/blogging-platform-api)

## Features

- **Create** blog posts with title, content, category, and tags
- **Read** all blog posts or retrieve individual posts by ID
- **Search** blog posts by title, content, or category
- **Update** existing blog posts while preserving creation timestamps
- **Delete** blog posts
- MongoDB integration for persistent storage
- JSON request/response format

## Getting Started

### Prerequisites

- Go 1.18 or higher
- MongoDB instance (cloud or local)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/dmandevv/blogging-platform-api.git
cd blogging-platform-api
```

2. Install dependencies:
```bash
go mod download
```

3. Configure MongoDB connection in `main.go` (update the connection string)

4. Run the application:
```bash
go run main.go
```

The API will start on `http://localhost:8080`

## API Endpoints

### Data Model

A blog post contains the following fields:
- `_id` (ObjectID): Unique identifier (auto-generated)
- `title` (string): Title of the blog post (required)
- `content` (string): Content/body of the blog post (required)
- `category` (string): Category of the blog post
- `tags` (array): Array of tags for the blog post
- `created_at` (timestamp): When the post was created (auto-set)
- `updated_at` (timestamp): When the post was last updated (auto-set)

### 1. Get All Blog Posts

**Endpoint:** `GET /posts`

**Description:** Retrieve all blog posts or search by terms

**Query Parameters:**
- `term` (optional): Search term to filter posts by title, content, or category

**Example - Get all posts:**
```bash
curl http://localhost:8080/posts
```

**Example - Search for posts:**
```bash
curl "http://localhost:8080/posts?term=golang"
```

**Response (200 OK):**
```json
[
  {
    "_id": "507f1f77bcf86cd799439011",
    "title": "Getting Started with Go",
    "content": "Go is a powerful programming language...",
    "category": "Programming",
    "tags": ["go", "programming", "tutorial"],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

### 2. Get Single Blog Post

**Endpoint:** `GET /posts/{_id}`

**Description:** Retrieve a specific blog post by ID

**Path Parameters:**
- `_id` (string): MongoDB ObjectID of the blog post

**Example:**
```bash
curl http://localhost:8080/posts/507f1f77bcf86cd799439011
```

**Response (200 OK):**
```json
{
  "_id": "507f1f77bcf86cd799439011",
  "title": "Getting Started with Go",
  "content": "Go is a powerful programming language...",
  "category": "Programming",
  "tags": ["go", "programming", "tutorial"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (404 Not Found):**
```json
"No document with this ID found: 507f1f77bcf86cd799439011"
```

### 3. Create Blog Post

**Endpoint:** `POST /posts`

**Description:** Create a new blog post

**Request Body:**
```json
{
  "title": "Getting Started with Go",
  "content": "Go is a powerful programming language...",
  "category": "Programming",
  "tags": ["go", "programming", "tutorial"]
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with Go",
    "content": "Go is a powerful programming language...",
    "category": "Programming",
    "tags": ["go", "programming", "tutorial"]
  }'
```

**Response (201 Created):**
```json
{
  "_id": "507f1f77bcf86cd799439011",
  "title": "Getting Started with Go",
  "content": "Go is a powerful programming language...",
  "category": "Programming",
  "tags": ["go", "programming", "tutorial"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (400 Bad Request):**
```
Title and Content are required
```

### 4. Update Blog Post

**Endpoint:** `PUT /posts/{_id}`

**Description:** Update an existing blog post. The creation timestamp is preserved, and the updated timestamp is automatically set.

**Path Parameters:**
- `_id` (string): MongoDB ObjectID of the blog post

**Request Body:**
```json
{
  "title": "Updated Title",
  "content": "Updated content here...",
  "category": "Programming",
  "tags": ["go", "updated"]
}
```

**Example:**
```bash
curl -X PUT http://localhost:8080/posts/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Advanced Go Patterns",
    "content": "Advanced patterns in Go programming...",
    "category": "Programming",
    "tags": ["go", "patterns", "advanced"]
  }'
```

**Response (201 Created):**
```json
{
  "_id": "507f1f77bcf86cd799439011",
  "title": "Advanced Go Patterns",
  "content": "Advanced patterns in Go programming...",
  "category": "Programming",
  "tags": ["go", "patterns", "advanced"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-16T14:22:15Z"
}
```

**Error Response (404 Not Found):**
```
Failed to find blog post in MongoDB
```

### 5. Delete Blog Post

**Endpoint:** `DELETE /posts/{_id}`

**Description:** Delete a blog post

**Path Parameters:**
- `_id` (string): MongoDB ObjectID of the blog post

**Example:**
```bash
curl -X DELETE http://localhost:8080/posts/507f1f77bcf86cd799439011
```

**Response (204 No Content):**
```
Blog deleted successfully
```

**Error Response (404 Not Found):**
```
Failed to find blog
```

## HTTP Status Codes

- `200 OK`: Successfully retrieved data
- `201 Created`: Resource successfully created or updated
- `204 No Content`: Resource successfully deleted
- `400 Bad Request`: Invalid request or missing required fields
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Database connection error

## Technologies Used

- **Language:** Go
- **Database:** MongoDB
- **HTTP:** Go standard library (`net/http`)
- **Driver:** MongoDB Go Driver v2

## License

See LICENSE file for details
