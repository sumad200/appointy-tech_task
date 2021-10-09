package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"username,omitempty" bson:"username,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"pwd,omitempty" bson:"pwd,omitempty"`
}

type Posts struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption    string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Imageurl   string             `json:"iurl,omitempty" bson:"iurl,omitempty"`
	PostedTime time.Time          `json:"ptime,omitempty" bson:"ptime,omitempty"`
}

func CreateUser(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response.Header().Set("content-type", "application/json")
		var person User
		_ = json.NewDecoder(request.Body).Decode(&person)
		s := person.Password
		h := sha256.New()
		h.Write([]byte(s))
		bs := h.Sum(nil)
		person.Password = base64.StdEncoding.EncodeToString(bs)
		collection := client.Database("api-data").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, _ := collection.InsertOne(ctx, person)
		json.NewEncoder(response).Encode(result)
	} else {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func CreatePost(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response.Header().Set("content-type", "application/json")
		var post Posts
		_ = json.NewDecoder(request.Body).Decode(&post)
		post.PostedTime = time.Now()
		collection := client.Database("api-data").Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, _ := collection.InsertOne(ctx, post)
		json.NewEncoder(response).Encode(result)
	} else {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		var people []User
		collection := client.Database("api-data").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var person User
			cursor.Decode(&person)
			people = append(people, person)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(people)
	} else {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func GetpostsEndpoint(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		var post []Posts
		collection := client.Database("api-data").Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var posts Posts
			cursor.Decode(&posts)
			post = append(post, posts)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(post)
	} else {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func GetUser(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		var getUserRe = regexp.MustCompile(`^\/users\/(.*)$`)
		matches := getUserRe.FindStringSubmatch(request.URL.Path)
		id, _ := primitive.ObjectIDFromHex(matches[1])
		var person User
		collection := client.Database("api-data").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := collection.FindOne(ctx, User{ID: id}).Decode(&person)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(person)
	} else {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func GetPost(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		var getUserRe = regexp.MustCompile(`^\/posts\/(.*)$`)
		matches := getUserRe.FindStringSubmatch(request.URL.Path)
		id, _ := primitive.ObjectIDFromHex(matches[1])
		var post Posts
		collection := client.Database("api-data").Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := collection.FindOne(ctx, Posts{ID: id}).Decode(&post)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(post)
	} else {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("Starting the API on port 0483")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	mux := http.NewServeMux()
	mux.HandleFunc("/users", CreateUser)
	mux.HandleFunc("/allusers", GetPeopleEndpoint)
	mux.HandleFunc("/allposts", GetpostsEndpoint)
	mux.HandleFunc("/users/", GetUser)
	mux.HandleFunc("/posts/", GetPost)
	mux.HandleFunc("/posts", CreatePost)
	http.ListenAndServe(":0483", mux)
}
