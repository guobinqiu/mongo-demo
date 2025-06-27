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

type Result struct {
	RoleName string `bson:"_id"`
	Count    int    `bson:"count"`
}

// 统计每种角色下有多少用户
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	_ = seed.SeedUsers(ctx, collection)

	// MongoDB 聚合操作 (类似SQL的 GROUP BY role.name)
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$role.name"},                      // 分组字段
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}, // 每组累加计数
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "count", Value: -1}, // 降序
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var r Result
		if err := cursor.Decode(&r); err != nil {
			panic(err)
		}
		fmt.Printf("角色: %s，有 %d 个用户\n", r.RoleName, r.Count)
	}
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// db.users.aggregate([
//   {
//     $group: {
//       _id: "$role.name",
//       count: { $sum: 1 }
//     }
//   },
//   {
//     $sort: { count: -1 }
//   }
// ])
