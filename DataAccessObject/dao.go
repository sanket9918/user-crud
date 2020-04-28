package dataaccessobject

import (
	"context"
	"log"

	"golang-rest-api-mongo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// DAO declaration
type DAO struct {
	Server   string
	Database string
}

var db *mongo.Database

// COLLECTION declaration
const (
	COLLECTION = "users"
)

// Connection to database
func (m *DAO) Connection() {
	clientOpts := options.Client().ApplyURI(m.Server)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(m.Database)
}

// FindAll list of users
func (m *DAO) FindAll() (users []models.User, err error) {
	cursor, err := db.Collection(COLLECTION).Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Fatal(err)
	}
	return users, err
}

// FindByID will find a user by its id
func (m *DAO) FindByID(id string) (user models.User, err error) {
	opts := options.FindOne().SetSort(bson.D{})
	err = db.Collection(COLLECTION).FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, opts).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return user, err
		}
		log.Fatal(err)
	}
	return user, err
}

// Insert a user into database
func (m *DAO) Insert(user models.User) (err error) {
	_, err = db.Collection(COLLECTION).InsertOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: &user}})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Delete an existing user
func (m *DAO) Delete(user models.User) (err error) {
	opts := options.FindOneAndDelete().SetProjection(bson.D{primitive.E{Key: "_id", Value: &user}})
	err = db.Collection(COLLECTION).FindOneAndDelete(context.TODO(), bson.D{primitive.E{Key: "_id", Value: &user}}, opts).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		log.Fatal(err)
	}
	return err
}

// Update an existing user
func (m *DAO) Update(user models.User) (err error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "_id", Value: &user}}}}
	_, err = db.Collection(COLLECTION).UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
