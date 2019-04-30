package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	databaseName = "go_database"
	collectionName = "user"
)

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`

}

type finalResult struct {
	Status bool
	Message string
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtKey = []byte("captainjacksparrow")
var dB *mongo.Database
var collection *mongo.Collection

func generateJWT(user User) (string, int64) {

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err!=nil {
		log.Fatal(err)
	}

	return tokenString, expirationTime.Unix()
}

func authorizeJWT(tokenString string) bool {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !token.Valid {
		return false
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false
		}
		return false
	}
	return true
}

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

			tokenString, expirationTime := generateJWT(user)
			_, _ = strconv.ParseInt(string(expirationTime), 10, 64)

			if err!=nil {
				log.Fatal(err)
			}

			w.Header().Set("token",tokenString)

			finalData.Status = true
			finalData.Message = "password correct"
			_ = json.NewEncoder(w).Encode(finalData)

		} else {
			finalData.Status = false
			finalData.Message = "password incorrect"
			_ = json.NewEncoder(w).Encode(finalData)
		}
	} else {
		finalData.Status = false
		finalData.Message = "user not registered"
		_ = json.NewEncoder(w).Encode(finalData)
	}
}


func deleteUser(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("token")

	decoder := json.NewDecoder(r.Body)
	var formData User
	err := decoder.Decode(&formData)
	if err!=nil {
		log.Fatal(err)
	}

	var user User
	user.Email = formData.Email
	user.Password = formData.Password

	authorize := authorizeJWT(tokenString)
	result := findUser(collection, user)
	var finalData finalResult
	if authorize {
		if len(result.Email) > 0 {
			if result.Password == formData.Password {
				if deleteUse(collection, user) {
					finalData.Status = true
					finalData.Message = "delete successful"
					_ = json.NewEncoder(w).Encode(finalData)
				} else {
					finalData.Status = false
					finalData.Message = "delete successful"
					_ = json.NewEncoder(w).Encode(finalData)
				}
			} else {
				finalData.Status = false
				finalData.Message = "password incorrect"
				_ = json.NewEncoder(w).Encode(finalData)
			}
		} else {
			finalData.Status = false
			finalData.Message = "user not registered"
			_ = json.NewEncoder(w).Encode(finalData)
		}
	} else {
		finalData.Status = false
		finalData.Message = "jwt not authorized"
		_ = json.NewEncoder(w).Encode(finalData)
	}

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
		localUser.Email = ""
		localUser.Password = ""
	}

	return localUser
}

func deleteUse(collection *mongo.Collection, user User) bool {

	filter := bson.D{{"email", user.Email}}

	result, err := collection.DeleteMany(context.TODO(), filter)
	if err!=nil {
		return false
	}
	if result.DeletedCount > 0 {
		return true
	} else {
		return false
	}
}

func main() {

	// Database Connections
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

	// Routers
	router := mux.NewRouter()
	router.HandleFunc("/register", registerUser).Methods("POST")
	router.HandleFunc("/login", loginUser).Methods("GET")
	router.HandleFunc("/delete", deleteUser).Methods("DELETE")

	// Server Listener
	err = http.ListenAndServe(":3000",router)
	if err!=nil {
		log.Fatal(err)
	}
}