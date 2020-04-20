package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
// User type 
// Represents a user, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type User struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Age  int `bson:"age" json:"age"`
	Email string `bson:"email" json:"email"`
}



