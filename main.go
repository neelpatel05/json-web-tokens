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
	"os"
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

	status bool
	message string
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

	result := createUser(collection, user)
	var finalData finalResult
	if result == true {
		finalData.status = true
		finalData.message = "registered successfully"
		_ = json.NewEncoder(w).Encode(finalData)
	} else {
		finalData.status = false
		finalData.message = "not registered"
		_ = json.NewEncoder(w).Encode(finalData)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {

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

func findUser(collection *mongo.Collection, user User) bool {

	var user User
	filter := bson.D{
		{
			"email",os.Getenv("EMAIL"),
		},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println(user.Email)
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
	router.HandleFunc("/resgister", registerUser).Methods("POST")
	router.HandleFunc("/login", loginUser).Methods("GET")
	router.HandleFunc("/logout", logoutUser).Methods("GET")
	router.HandleFunc("/delete", deleteUser).Methods("GET")

	err = http.ListenAndServe(":3000",router)
	if err!=nil {
		log.Fatal(err)
	}
}