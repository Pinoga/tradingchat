package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOptions struct {
	URI string
	DB  string
}

type MongoStore struct {
	Client   *mongo.Client
	Database string
}

func NewStore(op MongoOptions) (*MongoStore, error) {
	client, err := getClient(op)
	if err != nil {
		return nil, err
	}
	return &MongoStore{
		Client:   client,
		Database: op.DB,
	}, nil
}

func ToBson(v map[string]interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func getClient(op MongoOptions) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(op)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(op.URI))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *MongoStore) defaultDB() *mongo.Database {
	return s.Client.Database(s.Database)
}

func (s *MongoStore) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
		if err := s.Client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (s *MongoStore) InsertOne(collection string, entity interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := s.defaultDB().Collection(collection)

	r, err := c.InsertOne(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("error inserting document into %v: %s", collection, err)
	}

	if oid, ok := r.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	} else {
		return "", fmt.Errorf("error getting id string from %v", r.InsertedID)
	}
}

func (s *MongoStore) FindOne(collection string, filter map[string]interface{}, dest interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := s.defaultDB().Collection(collection)

	id, ok := filter["_id"]
	if ok {
		idString, ok := id.(string)
		err := fmt.Errorf("invalid _id paramenter value: %v", idString)
		if !ok {
			return false, err
		}
		objID, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			return false, err
		}
		filter["_id"] = objID
	}

	f, err := ToBson(filter)
	if err != nil {
		fmt.Printf("error converting %v filter to bson: %v", filter, err)
		return false, err
	}

	err = c.FindOne(ctx, *f).Decode(dest)
	if err != nil && err == mongo.ErrNoDocuments {
		return false, err
	} else if err != nil {
		fmt.Printf("FindOne error: %v", err)
		return false, err
	}
	return true, nil
}
