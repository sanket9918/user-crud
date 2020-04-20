package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"golang-rest-api-mongo/config"
	"golang-rest-api-mongo/dao"
	"golang-rest-api-mongo/models"
)

var config = Config{}
var dao = UsersDAO{}

// AllUsersEndPoint will GET list of users
func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

// FindUserEndpoint will GET a users by its ID
func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	if params := r.Context().Value(varsKey); params != nil {
		return params.(map[string]string)
	}
	user, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// CreateUserEndPoint will POST a new user
func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = bson.NewObjectId()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUserEndPoint will PUT update an existing user
func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if params := r.Context().Value(varsKey); params != nil {
		return params.(map[string]string)
	}
	var user User
	user.ID = bson.ObjectIdHex(params["id"])
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DeleteUserEndPoint will DELETE an existing user
func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	http.HandleFunc("/users", AllUsersEndPoint).Methods("GET")
	http.HandleFunc("/users", CreateUserEndPoint).Methods("POST")
	http.HandleFunc("/users/{id}", UpdateUserEndPoint).Methods("PUT")
	http.HandleFunc("/users", DeleteUserEndPoint).Methods("DELETE")
	http.HandleFunc("/users/{id}", FindUserEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
