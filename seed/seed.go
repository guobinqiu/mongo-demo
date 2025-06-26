package seed

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedUsers(ctx context.Context, collection *mongo.Collection) error {
	// 清空原始数据
	_ = collection.Drop(ctx)

	// 插入样例数据
	users := []any{
		bson.M{
			"username": "alice",
			"email":    "alice@example.com",
			"age":      25,
			"role": bson.M{
				"name":        "admin",
				"permissions": []string{"read", "write", "delete"},
			},
		},
		bson.M{
			"username": "bob",
			"email":    "bob@example.com",
			"age":      30,
			"role": bson.M{
				"name":        "user",
				"permissions": []string{"read"},
			},
		},
		bson.M{
			"username": "charlie",
			"email":    "charlie@example.com",
			"age":      28,
			"role": bson.M{
				"name":        "moderator",
				"permissions": []string{"read", "write"},
			},
		},
	}

	_, err := collection.InsertMany(ctx, users)
	return err
}
