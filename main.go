package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

type Comment struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Time    time.Time          `json:"time" bson:"time"`
	Comment string             `json:"comment" bson:"comment"`
}

var client *mongo.Client
var comments []Comment

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.3:27017") // mongodb://mongo:27017

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Route handles & endpoints
	router.HandleFunc("/comments", GetComments).Methods("GET")
	router.HandleFunc("/comments", CreateComment).Methods("POST")
	router.HandleFunc("/comments/{id}", GetComment).Methods("GET")
	router.HandleFunc("/comments/{id}", UpdateComment).Methods("PATCH")
	router.HandleFunc("/comments/{id}", DeleteComment).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetComments retrieves all comments
func GetComments(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(comments)
}

// GetComment retrieves a single comment by id
func GetComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range comments {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Comment{})
}

// CreateComment creates a new comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	_ = json.NewDecoder(r.Body).Decode(&comment)
	comment.ID = fmt.Sprintf("%d", len(comments)+1)
	comments = append(comments, comment)
	json.NewEncoder(w).Encode(comment)
}

// UpdateComment updates an existing comment by id
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var comment Comment
	_ = json.NewDecoder(r.Body).Decode(&comment)
	for i, item := range comments {
		if item.ID == params["id"] {
			comment.ID = params["id"]
			comments[i] = comment
			json.NewEncoder(w).Encode(comment)
			return
		}
	}
	json.NewEncoder(w).Encode(comments)
}

// DeleteComment deletes a comment by id
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, item := range comments {
		if item.ID == params["id"] {
			comments = append(comments[:i], comments[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(comments)
}
