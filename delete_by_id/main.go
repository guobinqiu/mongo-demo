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

	// 确保有数据
	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	// 查找一个用户，取出ID
	var user model.User
	err = collection.FindOne(ctx, bson.M{"username": "jack"}).Decode(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("准备删除用户: %s, ID: %s\n", user.Username, user.ID.Hex())

	// 根据ID删除
	filter := bson.M{"_id": user.ID}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("删除了 %d 条记录\n", result.DeletedCount)
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// var user = db.users.findOne({ "username": "jack" });
// if (user) {
// 	db.users.deleteOne({ "_id": user._id });
// }
