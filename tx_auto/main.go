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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0")) // 注意加 ?replicaSet=
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	// 清空集合（非事务操作）
	_ = collection.Drop(ctx)

	// 开始会话（开启事务需要会话）
	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	// 执行事务
	if _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (any, error) {
		user := model.User{
			Username: "charlie",
			Email:    "charlie@example.com",
			Age:      28,
			Role: model.Role{
				Name:        "admin",
				Permissions: []string{"read", "write", "delete"},
			},
		}
		result, err := collection.InsertOne(sessCtx, user)
		if err != nil {
			return nil, err // 返回错误，事务自动回滚，不用你手动调用rollback
		}

		fmt.Printf("事务中插入成功，ID: %v\n", result.InsertedID)
		return result.InsertedID, nil
	}); err != nil {
		panic(err)
	}
}
