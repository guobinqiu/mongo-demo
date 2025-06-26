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

	// 更新条件：所有年龄小于30的用户
	filter := bson.M{"age": bson.M{"$lt": 30}}

	// 更新内容：年龄 +1，角色权限替换
	update := bson.M{
		"$inc": bson.M{"age": 1}, // 年龄加1
		"$set": bson.M{
			"role.permissions": []string{"read", "update"},
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Printf("匹配到 %d 条，实际更新 %d 条\n", result.MatchedCount, result.ModifiedCount)
}

// mongosh mongodb://localhost:27017/testdb
// show collections
// db.users.updateMany(
//  { "age": {"$lt": 30} },
// 	{
//     $inc: { age: 1 },
//     $set: {
//       "role.permissions": ["read", "update"]
//     }
//   }
// )
