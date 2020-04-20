package dao

import (
	"log"

	"golang-rest-api-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)
// UsersDAO declaration
type UsersDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

// COLLECTION declaration
const (
	COLLECTION = "users"
)

// Connection to database
func (m *UsersDAO) Connection() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// FindAll list of users
func (m *UsersDAO) FindAll() ([]models.User, error) {
	var users []models.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

// FindByID will find a user by its id
func (m *UsersDAO) FindByID(id string) (models.User, error) {
	var user models.User
	err := db.C(COLLECTION).FindId(primitive.ObjectID(id)).One(&user)
	return user, err
}

// Insert a user into database
func (m *UsersDAO) Insert(user models.User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// Delete an existing user
func (m *UsersDAO) Delete(user models.User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

// Update an existing user
func (m *UsersDAO) Update(user models.User) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}
