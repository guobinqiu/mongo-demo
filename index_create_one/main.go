package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: 1}}, // 正向索引从小到大排序
		Options: options.Index().
			SetName("idx_username").
			SetUnique(true),
	}

	// Drop 索引（忽略不存在的报错）
	name := *indexModel.Options.Name
	_, _ = collection.Indexes().DropOne(ctx, name) // 不 panic，让它静默失败

	name, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic(err)
	}
	fmt.Printf("索引创建成功: %s\n", name)
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// db.users.dropIndex("idx_username")
// db.users.createIndex({ username: 1 }, { name: "idx_username", unique: true })
