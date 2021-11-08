package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOptions struct {
	URI string
	DB  string
}

type MongoStore struct {
	DB *mongo.Database
}

func NewStore(op MongoOptions) *MongoStore {
	db := getMongoClient(op)
	return &MongoStore{
		DB: db,
	}
}

func ToBson(v map[string]interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func getMongoClient(op MongoOptions) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(op.URI))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return client.Database(op.DB)
}

func (s *MongoStore) InsertOne(collection string, entity interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := s.DB.Collection(collection)

	r, err := c.InsertOne(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("error inserting document into %v: %s", collection, err)
	}

	fmt.Println(r.InsertedID)
	return r.InsertedID, nil
}

func (s *MongoStore) FindOne(collection string, filter map[string]interface{}, dest interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := s.DB.Collection(collection)

	f, err := ToBson(filter)
	if err != nil {
		fmt.Printf("error converting %v filter to bson: %v", filter, err)
		return false, err
	}

	err = c.FindOne(ctx, *f).Decode(dest)
	if err != nil && err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
