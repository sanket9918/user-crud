package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User type
// Represents a user, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type User struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Name  string             `bson:"name" json:"name"`
	Age   int                `bson:"age" json:"age"`
	Email string             `bson:"email" json:"email"`
}

// DAO declaration
type DAO struct {
	Server   string
	Database string
}

// Database variable declaration
var (
	db        *mongo.Database
	user      User
	err       error
	dao       = DAO{}
	validPath = regexp.MustCompile(`^/users/(update|delete|find)/([a-z0-9]+)$`)
)

// COLLECTION declaration
const (
	COLLECTION = "users"
)

// Connection to database
func (m *DAO) Connection() {
	clientOptions := options.Client().ApplyURI(m.Server)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	db = client.Database(m.Database)
}

// FindAll list of users
func (m *DAO) FindAll() (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "age", Value: -1}})
	cursor, err := db.Collection(COLLECTION).Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	return users, err
}

// FindByID will find a user by its id
func (m *DAO) FindByID(id string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Collection(COLLECTION).FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&user)
	defer cancel()
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return user, err
}

// Delete an existing user
func (m *DAO) Delete(id string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Collection(COLLECTION).FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return user, err
}

// Insert a user into database
func (m *DAO) Insert(user User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = db.Collection(COLLECTION).InsertOne(ctx, &user)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Update an existing user
func (m *DAO) Update(user User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: &user}}
	_, err = db.Collection(COLLECTION).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

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
	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUserEndpoint will DELETE an existing user
func DeleteUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
	deletedUser, err := dao.Delete(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, deletedUser)
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
		fmt.Println(r.URL.EscapedPath())
		m := validPath.FindStringSubmatch(r.URL.EscapedPath())
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[len(m)-1])
	}
}

// Parse the configuration file 'conf.json', and establish a connection to DB
func init() {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal("error:", err)
	}
	decoder := json.NewDecoder(file)
	defer file.Close()
	err = decoder.Decode(&dao)
	if err != nil {
		log.Fatal("error:", err)
	}

	dao.Connection()
}

// Define HTTP request routes
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", AllUsersEndpoint)
	mux.HandleFunc("/users/new", CreateUserEndpoint)
	mux.Handle("/users/update/", makeHandler(UpdateUserEndpoint))
	mux.Handle("/users/delete/", makeHandler(DeleteUserEndpoint))
	mux.Handle("/users/find/", makeHandler(FindUserEndpoint))
	if err = http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
