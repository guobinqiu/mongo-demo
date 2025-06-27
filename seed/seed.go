package seed

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedUsers(ctx context.Context, collection *mongo.Collection) error {
	// 清空原始数据
	_ = collection.Drop(ctx)

	// 插入样例数据（10条）
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
		bson.M{
			"username": "jack",
			"email":    "jack@example.com",
			"age":      42,
			"role": bson.M{
				"name":        "user",
				"permissions": []string{"read"},
			},
		},
		bson.M{
			"username": "emma",
			"email":    "emma@example.com",
			"age":      22,
			"role": bson.M{
				"name":        "user",
				"permissions": []string{"read"},
			},
		},
		bson.M{
			"username": "frank",
			"email":    "frank@example.com",
			"age":      35,
			"role": bson.M{
				"name":        "admin",
				"permissions": []string{"read", "write", "delete"},
			},
		},
		bson.M{
			"username": "grace",
			"email":    "grace@example.com",
			"age":      27,
			"role": bson.M{
				"name":        "editor",
				"permissions": []string{"read", "write"},
			},
		},
		bson.M{
			"username": "harry",
			"email":    "harry@example.com",
			"age":      33,
			"role": bson.M{
				"name":        "user",
				"permissions": []string{"read"},
			},
		},
		bson.M{
			"username": "ivy",
			"email":    "ivy@example.com",
			"age":      29,
			"role": bson.M{
				"name":        "moderator",
				"permissions": []string{"read", "write"},
			},
		},
		bson.M{
			"username": "leo",
			"email":    "leo@example.com",
			"age":      38,
			"role": bson.M{
				"name":        "editor",
				"permissions": []string{"read", "write"},
			},
		},
	}

	_, err := collection.InsertMany(ctx, users)
	return err
}
