package model

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role"`
	Hash     string `json:"hash" bson:"hash"`
}

func (u *User) GetID() string {
	return u.ID
}
