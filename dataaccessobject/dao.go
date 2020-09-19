package dataaccessobject

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BryanSouza91/user-crud/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DAO declaration
type DAO struct {
	Server   string
	Database string
}

// Database variable declaration
var db *mongo.Database

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
func (m *DAO) FindAll() (users []models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	results := bson.M{}
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "age", Value: -1}})
	cursor, err := db.Collection(COLLECTION).Find(ctx, results, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	return users, err
}

// FindByID will find a user by its id
func (m *DAO) FindByID(id string) (user models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	result := db.Collection(COLLECTION).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	err = result.Decode(&user)
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

// Insert a user into database
func (m *DAO) Insert(user models.User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = db.Collection(COLLECTION).InsertOne(ctx, &user)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Delete an existing user
func (m *DAO) Delete(id string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	opts := options.FindOneAndDelete().SetProjection(bson.D{primitive.E{Key: "_id", Value: id}})
	err = db.Collection(COLLECTION).FindOneAndDelete(ctx, bson.D{primitive.E{Key: "_id", Value: id}}, opts).Decode(&id)
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
