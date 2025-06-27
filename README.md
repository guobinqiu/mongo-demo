# Mongo Demo

本项目用于演示 MongoDB 在 Go 中的常见操作，包括文档的增删改查（CRUD）、分页查询、索引管理、聚合操作以及事务处理等。

每个操作示例单独放在对应目录，便于学习和测试。

## 安装

```
./setup.sh
```

为了支持事务功能，MongoDB 以副本集模式启动

## 功能模块列表

| 目录                | 功能描述         |
| ------------------- | ---------------- |
| aggregate/          | 聚合查询示例     |
| count/              | 文档计数示例     |
| delete_by_id/       | 根据ID删除文档   |
| delete_many/        | 删除多个文档     |
| delete_one/         | 删除单个文档     |
| find_all/           | 查询所有文档     |
| find_by_like/       | 模糊查询示例     |
| find_by_role_field/ | 根据角色字段查询 |
| find_by_user_field/ | 根据用户字段查询 |
| find_one/           | 查询单个文档     |
| find_pagination/    | 分页查询示例     |
| index_create_many/  | 创建多个索引     |
| index_create_one/   | 创建单个索引     |
| insert_many/        | 插入多个文档     |
| insert_one/         | 插入单个文档     |
| tx_auto/            | 自动事务示例     |
| tx_manual/          | 手动事务示例     |
| update_by_id/       | 根据ID更新文档   |
| update_many/        | 更新多个文档     |
| update_one/         | 更新单个文档     |

## MongoDB Shell 命令示例

| 功能         | Shell 命令                                                                                     |
| ------------ | ---------------------------------------------------------------------------------------------- |
| 查询所有文档 | `db.users.find()`                                                                              |
| 分页查询     | `db.users.find().skip(0).limit(10)`                                                            |
| 根据ID查询   | `db.users.findOne({_id: ObjectId("id")})`                                                      |
| 根据字段查询 | `db.users.find({age: {$gt: 25}})`                                                              |
| 模糊查询     | `db.users.find({name: /pattern/})`                                                             |
| 插入文档     | `db.users.insertOne({name: "John", age: 30})`                                                  |
| 更新文档     | `db.users.updateOne({_id: ObjectId("id")}, {$set: {age: 31}})`                                 |
| 删除文档     | `db.users.deleteOne({_id: ObjectId("id")})`                                                    |
| 聚合查询     | `db.users.aggregate([{$match: {age: {$gt: 25}}}, {$group: {_id: "$role", count: {$sum: 1}}}])` |
| 创建索引     | `db.users.createIndex({name: 1})`                                                              |
| 自动事务     | `session.withTransaction(() => {                                                               |
  db.users.insertOne({...}, {session});
  return true;
});` |
| 手动事务 | `session.startTransaction();
try {
  db.users.insertOne({...}, {session});
  session.commitTransaction();
} catch(e) {
  session.abortTransaction();
}` |