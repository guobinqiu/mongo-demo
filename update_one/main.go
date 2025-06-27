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

	// 更新条件：用户名为 "bob"
	filter := bson.M{"username": "bob"}

	// 更新内容：设置年龄为 32，角色权限增加 "update"
	update := bson.M{
		"$set": bson.M{
			"age":              32,
			"role.permissions": []string{"read", "update"}, // 替换权限数组
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Printf("匹配到 %d 条，实际更新 %d 条\n", result.MatchedCount, result.ModifiedCount)
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// db.users.updateOne(
//   { username: "bob" },
//   { $set: { age: 32, "role.permissions": ["read", "update"] } }
// )
