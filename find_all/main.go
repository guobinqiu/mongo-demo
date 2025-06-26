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

	// 查询所有用户
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	fmt.Println("所有用户：")
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			panic(err)
		}
		fmt.Printf(" - %+v\n", user)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}
}

// mongosh mongodb://localhost:27017/testdb
// show collections
// db.users.find().pretty()
