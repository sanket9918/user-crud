package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"user-crud/config"
	"user-crud/dataaccessobject"
	"user-crud/models"
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
func FindUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
	user, err := dao.FindByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// CreateUserEndPoint will POST a new user
func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	user.ID = primitive.NewObjectID()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUserEndPoint will PUT update an existing user
func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request, id string) {
	var user models.User
	var err error
	user.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := dao.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DeleteUserEndPoint will DELETE an existing user
func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request, id string) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
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

var validPath = regexp.MustCompile("^/user/(update|delete|find)/([a-zA-Z0-9]+)$")

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

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	conf.Read()

	dao.Server = conf.Server
	dao.Database = conf.Database
	dao.Connection()
}

// Define HTTP request routes
func main() {
	http.HandleFunc("/users", AllUsersEndPoint)
	http.HandleFunc("/users/new", CreateUserEndPoint)
	http.HandleFunc("/users/update/{id}", makeHandler(UpdateUserEndPoint))
	http.HandleFunc("/users/delete/{id}", makeHandler(DeleteUserEndPoint))
	http.HandleFunc("/users/find/{id}", makeHandler(FindUserEndpoint))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
