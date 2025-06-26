package main

import (
	"context"
	"fmt"
	"time"

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

	// 确保有数据
	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	// 删除条件：年龄等于25的用户
	filter := bson.M{"age": bson.M{"$eq": 25}}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("删除了 %d 条记录\n", result.DeletedCount)
}

// mongosh mongodb://localhost:27017/testdb
// show collections
// db.users.deleteOne({"age": {"$eq": 25}})
