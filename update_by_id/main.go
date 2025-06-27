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

	// 确保数据存在
	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	// 先查一个用户，取ID
	var user model.User
	err = collection.FindOne(ctx, bson.M{"username": "bob"}).Decode(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("准备更新用户: %s, ID: %s\n", user.Username, user.ID.Hex())

	// 更新条件：_id
	filter := bson.M{"_id": user.ID}

	// 更新内容
	update := bson.M{
		"$set": bson.M{
			"age":              35,
			"role.permissions": []string{"read", "write", "execute"},
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Printf("匹配到 %d 条，实际更新 %d 条\n", result.MatchedCount, result.ModifiedCount)
}
