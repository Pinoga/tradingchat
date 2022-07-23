package store

import "tradingchat/pkg/mongodb"

type Store interface {
	InsertOne(string, mongodb.MongoEntity) (string, error)
	FindOne(string, map[string]interface{}, interface{}) error
	Disconnect()
}
