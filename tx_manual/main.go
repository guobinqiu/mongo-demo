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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

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
			_ = session.AbortTransaction(sessCtx)
			return err
		}
		fmt.Println("插入成功，ID:", result.InsertedID)

		if err := session.CommitTransaction(sessCtx); err != nil {
			return err
		}

		fmt.Println("事务提交成功")
		return nil
	})
	if err != nil {
		panic(err)
	}
}
