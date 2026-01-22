package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmandevv/blogging-platform-api/internal/config"
)

var tests = []struct {
	request          string
	httpStatus       int
	expectedResponse string
}{
	{
		request:          `{"title": "My First Blog Post", "content": "This is the content of my first blog post.", "category": "Technology", "tags": ["Tech", "Programming"]}`,
		httpStatus:       http.StatusOK,
		expectedResponse: "Blog created: My First Blog Post",
	},
	{
		request:          `{"title": "", "content": "This is the content of my first blog post.", "category": "Technology", "tags": ["Tech", "Programming"]}`,
		httpStatus:       http.StatusBadRequest,
		expectedResponse: "Title and Content are required\n",
	},
	{
		request:          `{"title": "My First Blog Post", "content": "", "category": "Technology", "tags": ["Tech", "Programming"]}`,
		httpStatus:       http.StatusBadRequest,
		expectedResponse: "Title and Content are required\n",
	},
	{
		request:          `{"name": "My First Blog Post", "content": "This is the content of my first blog post.", "category": "Technology", "tags": ["Tech", "Programming"]}`,
		httpStatus:       http.StatusBadRequest,
		expectedResponse: "Title and Content are required\n",
	},
}

func TestHandleCreate(t *testing.T) {

	for _, test := range tests {
		// Create a mock request with a POST method and a JSON body
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer([]byte(test.request)))
		req.Header.Set("Content-Type", "application/json")

		// Create a recorder to capture the response
		rr := httptest.NewRecorder()

		// Call the handler directly, passing the mock recorder and request
		HandleCreate(&config.Config{}, rr, req)

		// Check the status code
		if status := rr.Code; status != test.httpStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, test.httpStatus)
		}

		// Check the response body
		if rr.Body.String() != test.expectedResponse {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), test.expectedResponse)
		}
	}
}
