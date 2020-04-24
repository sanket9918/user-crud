package main

import (
	// "context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"

	"golang-rest-api-mongo/config"
	"golang-rest-api-mongo/dataaccessobject"
	"golang-rest-api-mongo/models"
)

var conf = config.Config{}
var dao = dataaccessobject.DAO{}

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
// func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
// 	var params interface{}
// 	if params := r.Context(); params != nil {
// 		return params
// 	}
// 	user, err := dao.FindByID(bson.D{primitive.E{Key: "_id", Value: params["id"]}})
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
// 		return
// 	}
// 	respondWithJSON(w, http.StatusOK, user)
// }

// CreateUserEndPoint will POST a new user
func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = primitive.NewObjectID()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUserEndPoint will PUT update an existing user
// func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	if params := r.Context(); params != nil {
// 		return params.(context)
// 	}
// 	var user models.User
// 	user.ID = primitive.NewObjectId(params["id"])
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	if err := dao.Update(user); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
// }

// DeleteUserEndPoint will DELETE an existing user
func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
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
	conf.Read()

	dao.Server = conf.Server
	dao.Database = conf.Database
	dao.Connection()
}

var validPath = regexp.MustCompile("^/users/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// Define HTTP request routes
func main() {
	http.HandleFunc("/users", makeHandler(AllUsersEndPoint))
	http.HandleFunc("/users", makeHandler(CreateUserEndPoint))
	http.HandleFunc("/users/{id}", makeHandler(UpdateUserEndPoint))
	http.HandleFunc("/users", makeHandler(DeleteUserEndPoint))
	http.HandleFunc("/users/{id}", makeHandler(FindUserEndpoint))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
