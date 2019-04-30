package main

import (
	"context"
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
	collectioName = "user"
)

type User struct {

	Email string
	Password string
}

func registerUser(w http.ResponseWriter, r *http.Request) {

}

func loginUser(w http.ResponseWriter, r *http.Request) {

}

func logoutUser(w http.ResponseWriter, r *http.Request) {

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

}

func createUser(collection *mongo.Collection) {

	var user User
	user.Email = os.Getenv("EMAIL")
	user.Password = os.Getenv("PASSWORD")

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult.InsertedID)
}

func findUser(collection *mongo.Collection) {

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

	os.Setenv("EMAIL", "dummy@gmail.com")
	os.Setenv("PASSWORD", "dummy")
	defer os.Unsetenv("EMAIL")
	defer os.Unsetenv("PASSWORD")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err!=nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err!=nil {
		log.Fatal(err)
	}

	var dB *mongo.Database
	var collection *mongo.Collection
	dB = client.Database(databaseName)
	collection = dB.Collection(collectioName)

	createUser(collection)
	findUser(collection)

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