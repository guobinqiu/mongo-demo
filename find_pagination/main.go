package main

import (
	"context"
	"fmt"
	"time"

	"github.com/guobinqiu/mongo-demo/model"
	"github.com/guobinqiu/mongo-demo/seed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PageResult 定义分页结构
type PageResult struct {
	Page       int64        // 当前页码
	PageSize   int64        // 每页条数
	TotalCount int64        // 总条数
	TotalPages int64        // 总页数
	Items      []model.User // 当前页数据
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("users")

	// 确保有数据
	if err := seed.SeedUsers(ctx, collection); err != nil {
		panic(err)
	}

	page := int64(1)     // 当前第几页（从1开始）
	pageSize := int64(3) // 每页几条
	filter := bson.M{}   // 查询条件

	result, err := FindUsersByPage(ctx, collection, page, pageSize, filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("总记录数: %d, 总页数: %d, 当前页: %d\n", result.TotalCount, result.TotalPages, result.Page)
	for _, user := range result.Items {
		fmt.Printf("- %+v\n", user)
	}
}

// FindUsersByPage 查询分页数据
func FindUsersByPage(ctx context.Context, collection *mongo.Collection, page, pageSize int64, filter bson.M) (*PageResult, error) {
	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 计算总页数
	// totalPages := int64(math.Ceil(float64(total) / float64(pageSize)))
	totalPages := (total + pageSize - 1) / pageSize

	// 跳过前 (page - 1) * pageSize 条数据，然后从第 n+1 条开始读取
	// skip 是从 0 开始计数的
	skip := (page - 1) * pageSize

	// 查询当前页数据
	opts := options.Find().
		SetSkip(skip).
		SetLimit(pageSize).
		SetSort(bson.D{{Key: "username", Value: 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	for cursor.Next(ctx) {
		var u model.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return &PageResult{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
		Items:      users,
	}, nil
}
