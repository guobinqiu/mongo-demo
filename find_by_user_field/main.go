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
	// 设置 MongoDB 连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	// 如果集合为空则插入样例数据
	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	// 查询条件：年龄大于 25
	filter := bson.M{
		"age": bson.M{"$gt": 25},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	fmt.Printf("查询条件: %v\n", filter)
	fmt.Println("匹配结果：")
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

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// db.users.find({ "age": { $gt: 25 } }).pretty()
