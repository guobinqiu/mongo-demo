package main

import (
	"context"
	"fmt"
	"time"

	"github.com/guobinqiu/mongo-demo/model"
	"github.com/guobinqiu/mongo-demo/seed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	// 查询条件：用户名为 "alice"
	filter := bson.M{"username": "alice"}

	var user model.User
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		panic(err)
	}

	fmt.Printf("查到用户：\n")
	fmt.Printf("ID: %s\nUsername: %s\nEmail: %s\nAge: %d\nRole: %s\nPermissions: %v\n",
		user.ID.Hex(), user.Username, user.Email, user.Age, user.Role.Name, user.Role.Permissions)
}

// mongosh mongodb://localhost:27017/testdb
// show collections
// db.users.findOne({"username": "alice"})
