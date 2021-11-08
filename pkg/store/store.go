package store

type Store interface {
	InsertOne(string, interface{}, interface{}) (bool, error)
	FindOne(string, map[string]interface{}, interface{}) (bool, error)
}
