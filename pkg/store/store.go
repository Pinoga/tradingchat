package store

type Store interface {
	InsertOne(string, interface{}) (string, error)
	FindOne(string, map[string]interface{}, interface{}) (bool, error)
	Disconnect()
}
