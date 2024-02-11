<br />
<div align="center">
  <h3 align="center">kyte</h3>

  <p align="center">
    MongoDB Query Builder for Golang
    <br />
    <br />
  </p>
</div>

<br>
<br>
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#motivation">Motivation</a></li>
    <li><a href="#installation">Installation</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#advanced-usage">Advanced Usage</a></li>
    <li><a href="#supported-operators">Supported Operators</a></li>
  </ol>
</details>
<br>
<br>

## About The Project

Kyte is a MongoDB query builder for Golang. It is designed to be simple and easy to use. Kyte's most unique feature is its ability to build MongoDB queries using a struct schema. This allows you to build queries using a struct that represents the schema of your MongoDB collection that prevents you from making mistakes in your queries check [Advanced Usage](#advanced-usage) for more details.

> Currently, it supports Filter operations only. Aggregate and Update operations will be added in the future.

## Motivation

Kyte is built with a straightforward vision: Make it easy. As any Gopher would know, working with MongoDB in Golang can be a bit of a hassle because of `bson` and `primitive` packages. Kyte aims to simplify the process of building MongoDB queries in Golang by providing a simple and easy-to-use API. It also provides a way to build queries using a struct schema, which is unique to Kyte.

## Installation

```sh
go get github.com/aaydin-tr/kyte
```

## Usage

### Basic Usage

Following is a simple example of how to use Kyte to build a MongoDB query.

```go
query, err := kyte.Filter().
    Equal("name", "John").
    GreaterThan("age", 20).
    Build()
```

The above code will generate the following MongoDB query:

```json
{ "name": {"$eq": "John"}, "age": {"$gt": 20 } }
```

### Advanced Usage

Kyte allows you to build queries using a struct schema. This is useful when you want to prevent mistakes in your queries. For example, if you have a struct that represents the schema of your MongoDB collection, you can use it to build queries. Kyte will ensure that the fields you use in your queries are valid and exist in the schema. For this purpose Kyte provides two option functions: `Source` and `ValidateField`.

`Source` function is used to specify the struct schema that you want to use to build queries. 
`ValidateField` function is used to validate the fields that you use in your queries. ( By default, it is enabled. )

> Note: You can also use `ValidateField(false)` to disable field validation.

```go

type User struct {
    Name string `bson:"name"`
    Age  int    `bson:"age"`
}

var user User

query, err := kyte.Filter(kyte.Source(&user)).
    Equal(&user.Name, "John").
    GreaterThan(&user.Age, 20).
    Build()
```

When you use `kyte.Source` function, you can pass a pointer of the struct field of your schema to the query builder functions. This will ensure that the fields you use in your queries are valid and exist in the schema.

> Note: You can also use `string` value as a field name and Kyte still will validate the field. *But using a pointer to the struct field is recommended.*

## Supported Operators

- Equal ([$eq](https://www.mongodb.com/docs/manual/reference/operator/query/eq/#mongodb-query-op.-eq))
  ```go
  Equal("name", "John")
  // { "name": {"$eq": "John"} }
  ```
- NotEqual ([$ne](https://www.mongodb.com/docs/manual/reference/operator/query/ne/#mongodb-query-op.-ne))
  ```go
  NotEqual("name", "John")
  // { "name": {"$ne": "John"} }
  ```
- GreaterThan ([$gt](https://www.mongodb.com/docs/manual/reference/operator/query/gt/#mongodb-query-op.-gt))
  ```go
  GreaterThan("age", 20)
  // { "age": {"$gt": 20} }
  ```
- GreaterThanOrEqual ([$gte](https://www.mongodb.com/docs/manual/reference/operator/query/gte/#mongodb-query-op.-gte))
  ```go
  GreaterThanOrEqual("age", 20)
  // { "age": {"$gte": 20} }
  ```
- LessThan ([$lt](https://www.mongodb.com/docs/manual/reference/operator/query/lt/#mongodb-query-op.-lt))
  ```go
  LessThan("age", 20)
  // { "age": {"$lt": 20} }
  ```
- LessThanOrEqual ([$lte](https://www.mongodb.com/docs/manual/reference/operator/query/lte/#mongodb-query-op.-lte))
  ```go
  LessThanOrEqual("age", 20)
  // { "age": {"$lte": 20} }
  ```
- In ([$in](https://www.mongodb.com/docs/manual/reference/operator/query/in/#mongodb-query-op.-in))
  ```go
  In("name", []string{"John", "Doe"})
  // { "name": {"$in": ["John", "Doe"]} }
  ```
- NotIn ([$nin](https://www.mongodb.com/docs/manual/reference/operator/query/nin/#mongodb-query-op.-nin))
  ```go
  NotIn("name", []string{"John", "Doe"})
  // { "name": {"$nin": ["John", "Doe"]} }
  ```
- And ([$and](https://www.mongodb.com/docs/manual/reference/operator/query/and/#mongodb-query-op.-and))
  ```go
  And(
      Filter().
          Equal("name", "John").
          GreaterThan("age", 20),
  )
  // { "$and": [ { "name": {"$eq": "John"} }, { "age": {"$gt": 20} } ] }
  ```
- Or ([$or](https://www.mongodb.com/docs/manual/reference/operator/query/or/#mongodb-query-op.-or))
  ```go
  Or(
      Filter().
          Equal("name", "John").
          GreaterThan("age", 20),
  )
  // { "$or": [ { "name": {"$eq": "John"} }, { "age": {"$gt": 20} } ] }
  ```
- Nor ([$nor](https://www.mongodb.com/docs/manual/reference/operator/query/nor/#mongodb-query-op.-nor))
  ```go
  Nor(
      Filter().
          Equal("name", "John").
          GreaterThan("age", 20),
  )
  // { "$nor": [ { "name": {"$eq": "John"} }, { "age": {"$gt": 20} } ] }
  ```
- Regex ([$regex](https://www.mongodb.com/docs/manual/reference/operator/query/regex/#mongodb-query-op.-regex))
  ```go
  Regex("name", regexp.MustCompile("John"), "i")
  // { "name": {"$regex": "John", "$options": "i"} }
  ```
- Exists ([$exists](https://www.mongodb.com/docs/manual/reference/operator/query/exists/#mongodb-query-op.-exists))
  ```go
  Exists("name", true)
  // { "name": {"$exists": true} }
  ```
- Type ([$type](https://www.mongodb.com/docs/manual/reference/operator/query/type/#mongodb-query-op.-type))
  ```go
  Type("name", bsontype.String)
  // { "name": {"$type": "string"} }
  ```
- Mod ([$mod](https://www.mongodb.com/docs/manual/reference/operator/query/mod/#mongodb-query-op.-mod) )
  ```go
  Mod("age", 2, 0)
  // { "age": {"$mod": [2, 0]} }
  ```
- Where ([$where](https://www.mongodb.com/docs/manual/reference/operator/query/where/#mongodb-query-op.-where))
  ```go
  Where("this.name.length > 10")
  // { "$where": "this.name.length > 10" }
  ```
- All ([$all](https://www.mongodb.com/docs/manual/reference/operator/query/all/#mongodb-query-op.-all) )
  ```go
  All("name", []string{"John", "Doe"})
  // { "name": {"$all": ["John", "Doe"]} }
  ```
- Size ([$size](https://www.mongodb.com/docs/manual/reference/operator/query/size/#mongodb-query-op.-size) )
  ```go
  Size("name", 10)
  // { "name": {"$size": 10} }
  ```
- JsonSchema ([$jsonSchema](https://www.mongodb.com/docs/manual/reference/operator/query/jsonSchema/#mongodb-query-op.-jsonSchema) )
  ```go
  JsonSchema(bson.M{"required": []string{"name"}})
  // { "$jsonSchema": {"required": ["name"]} }
  ```
- Raw
  ```go
  Raw(bson.D{{"name", "John"}})
  // { "name": "John" }
  ```




