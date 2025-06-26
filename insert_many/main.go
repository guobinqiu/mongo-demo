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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	// 清空原始数据
	_ = collection.Drop(ctx)

	// 创建多条用户数据
	users := []interface{}{
		model.User{
			Username: "david",
			Email:    "david@example.com",
			Age:      22,
			Role: model.Role{
				Name:        "user",
				Permissions: []string{"read"},
			},
		},
		model.User{
			Username: "eva",
			Email:    "eva@example.com",
			Age:      35,
			Role: model.Role{
				Name:        "admin",
				Permissions: []string{"read", "write", "delete"},
			},
		},
		model.User{
			Username: "frank",
			Email:    "frank@example.com",
			Age:      29,
			Role: model.Role{
				Name:        "moderator",
				Permissions: []string{"read", "write"},
			},
		},
	}

	// 批量插入
	result, err := collection.InsertMany(ctx, users)
	if err != nil {
		panic(err)
	}

	fmt.Println("成功插入的用户 ID 列表:")
	for _, id := range result.InsertedIDs {
		fmt.Printf(" - %v\n", id)
	}
}

// mongosh mongodb://localhost:27017/testdb
// show collections
// db.users.drop()
// db.users.insertMany([
//   {
//     "username": "david",
//     "email": "david@example.com",
//     "age": 22,
//     "role": {
//       "name": "user",
//       "permissions": ["read"]
//     }
//   },
//   {
//     "username": "eva",
//     "email": "eva@example.com",
//     "age": 35,
//     "role": {
//       "name": "admin",
//       "permissions": ["read", "write", "delete"]
//     }
//   },
//   {
//     "username": "frank",
//     "email": "frank@example.com",
//     "age": 29,
//     "role": {
//       "name": "moderator",
//       "permissions": ["read", "write"]
//     }
//   }
// ])
