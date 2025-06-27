package main

import (
	"context"
	"fmt"
	"time"

	"github.com/guobinqiu/mongo-demo/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 设置 MongoDB 连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0/?replicaSet=rs0"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	// 清空原始数据
	_ = collection.Drop(ctx)

	// 创建一条用户数据
	user := model.User{
		Username: "charlie",
		Email:    "charlie@example.com",
		Age:      28,
		Role: model.Role{
			Name:        "admin",
			Permissions: []string{"read", "write", "delete"},
		},
	}

	// 插入
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("成功插入用户，ID: %v\n", result.InsertedID)
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// db.users.drop()
// db.users.insertOne({
//   "username": "charlie",
//   "email": "charlie@example.com",
//   "age": 28,
//   "role": {
//     "name": "admin",
//     "permissions": ["read", "write", "delete"]
//   }
// })
