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

	// 定义多个索引
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "username", Value: 1}}, // username 升序索引
			Options: options.Index().
				SetName("idx_username").
				SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}}, // email 升序索引
			Options: options.Index().
				SetName("idx_email").
				SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "role.name", Value: 1}}, // 角色名索引
			Options: options.Index().
				SetName("idx_role_name"),
		},
	}

	// Drop 每个索引（忽略不存在的报错）
	for _, index := range indexModels {
		name := *index.Options.Name
		_, _ = collection.Indexes().DropOne(ctx, name) // 不 panic，让它静默失败
	}

	name, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}
	fmt.Printf("索引创建成功: %s\n", name)
}

// mongosh mongodb://localhost:27017/?replicaSet=rs0/testdb
// show collections
// ["idx_username", "idx_email", "idx_role_name"].forEach(function(name) {
//     try {
//         db.users.dropIndex(name);
//         print("索引 " + name + " 已删除");
//     } catch (e) {
//         print("未找到索引 " + name + "，跳过删除");
//     }
// });

// var indexModels = [
//     {
//         key: { username: 1 },
//         name: "idx_username",
//         unique: true
//     },
//     {
//         key: { email: 1 },
//         name: "idx_email",
//         unique: true
//     },
//     {
//         key: { "role.name": 1 },
//         name: "idx_role_name"
//     }
// ];

// var result = db.runCommand({
//     createIndexes: "users",
//     indexes: indexModels
// });

// if (result.ok === 1) {
//     print("所有索引创建成功！");
//     printjson(result);
// } else {
//     print("索引创建失败：");
//     printjson(result);
// }
