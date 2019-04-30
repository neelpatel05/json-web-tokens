package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const (
	databaseName = "go_database"
	collectionName = "user"
)

type User struct {
	Email string
	Password string
}

type finalResult struct {
	Status bool
	Message string
}

var dB *mongo.Database
var collection *mongo.Collection

func registerUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var formData User
	err := decoder.Decode(&formData)
	if err!=nil {
		log.Fatal(err)
	}

	var user User
	user.Email = formData.Email
	user.Password = formData.Password

	result1 := findUser(collection, user)
	var finalData finalResult
	fmt.Println(len(result1.Email))
	if len(result1.Email) == 0 {
		result2 := createUser(collection, user)
		if result2 == true {
			finalData.Status = true
			finalData.Message = "registered successfully"
			_ = json.NewEncoder(w).Encode(finalData)
		} else {
			finalData.Status = false
			finalData.Message = "not registered"
			_ = json.NewEncoder(w).Encode(finalData)
		}
	} else {
		finalData.Status = false
		finalData.Message = "user already exists"
		_ = json.NewEncoder(w).Encode(finalData)
	}

}

func loginUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var formData User
	err := decoder.Decode(&formData)
	if err!=nil {
		log.Fatal(err)
	}

	var user User
	user.Email = formData.Email
	user.Password = formData.Password

	result := findUser(collection, user)
	var finalData finalResult
	if len(result.Email)>0 {
		if result.Password == formData.Password {
			finalData.Status = true
			finalData.Message = "authorized"
			_ = json.NewEncoder(w).Encode(finalData)
		} else {
			finalData.Status = false
			finalData.Message = "not authorized"
			_ = json.NewEncoder(w).Encode(finalData)
		}
	} else {
		finalData.Status = false
		finalData.Message = "user not registered"
		_ = json.NewEncoder(w).Encode(finalData)
	}
}

func logoutUser(w http.ResponseWriter, r *http.Request) {

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

}

func createUser(collection *mongo.Collection, user User) bool {

	_, err := collection.InsertOne(context.TODO(), user)
	if err!=nil {
		return false
	} else {
		return true
	}

}

func findUser(collection *mongo.Collection, user User) User {

	var localUser User
	filter := bson.D{{"email", user.Email}}

	err := collection.FindOne(context.TODO(), filter).Decode(&localUser)
	if err!=nil {
		fmt.Println(localUser)
	}

	return localUser
}

func main() {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err!=nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err!=nil {
		log.Fatal(err)
	}

	dB = client.Database(databaseName)
	collection = dB.Collection(collectionName)

	router := mux.NewRouter()
	router.HandleFunc("/register", registerUser).Methods("POST")
	router.HandleFunc("/login", loginUser).Methods("GET")
	router.HandleFunc("/logout", logoutUser).Methods("GET")
	router.HandleFunc("/delete", deleteUser).Methods("GET")

	err = http.ListenAndServe(":3000",router)
	if err!=nil {
		log.Fatal(err)
	}
}