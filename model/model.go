package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	Name        string   `bson:"name" json:"name"`
	Permissions []string `bson:"permissions" json:"permissions"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Age      int                `bson:"age" json:"age"`
	Role     Role               `bson:"role" json:"role"`
}
