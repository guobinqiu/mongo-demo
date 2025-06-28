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

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	// 自动插入样例数据（如果为空）
	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	// 查询用户名包含 "a" 的用户（相当于 SQL: WHERE username LIKE '%a%'）
	filter := bson.M{
		// "username": primitive.Regex{Pattern: "a", Options: "i"}, // i = 忽略大小写
		"username": bson.M{"$regex": "a", "$options": "i"},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	fmt.Printf("查询条件: %v\n", filter)
	fmt.Println("匹配用户名包含 'a' 的用户：")
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			panic(err)
		}
		fmt.Printf(" - %s (%s)\n", user.Username, user.Email)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// db.users.find({
//   "username": {
//     $regex: "a",
//     $options: "i"
//   }
// }).pretty()
