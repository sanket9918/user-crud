package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/BryanSouza91/user-crud/dataaccessobject"
	"github.com/BryanSouza91/user-crud/models"
)

var (
	dao       = dataaccessobject.DAO{}
	validPath = regexp.MustCompile(`^/users/(update|delete|find)/([a-z0-9]+)$`)
)

// AllUsersEndpoint will GET list of users
func AllUsersEndpoint(w http.ResponseWriter, r *http.Request) {
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

// CreateUserEndpoint will POST a new user
func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
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

// UpdateUserEndpoint will PUT update an existing user
func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
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

// DeleteUserEndpoint will DELETE an existing user
func DeleteUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
	if err := dao.Delete(id); err != nil {
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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		m := validPath.FindStringSubmatch(r.URL.EscapedPath())
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// Parse the configuration file 'conf.json', and establish a connection to DB
func init() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&dao)
	if err != nil {
		log.Fatal("error:", err)
	}

	dao.Connection()
}

// Define HTTP request routes
func main() {
	http.HandleFunc("/users", AllUsersEndpoint)
	http.HandleFunc("/users/new", CreateUserEndpoint)
	http.HandleFunc("/users/update/", makeHandler(UpdateUserEndpoint))
	http.HandleFunc("/users/delete/", makeHandler(DeleteUserEndpoint))
	http.HandleFunc("/users/find/", makeHandler(FindUserEndpoint))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
