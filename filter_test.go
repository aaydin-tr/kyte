package kyte

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	t.Run("without options", func(t *testing.T) {
		filter := Filter()
		if filter == nil {
			t.Error("Filter should not be nil")
		}

		if filter.kyte.source != nil {
			t.Error("Filter.kyte should be nil")
		}

		if filter.kyte.checkField != false {
			t.Error("Filter.kyte.checkField should be false")
		}
	})

	t.Run("with options", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}

		filter := Filter(Source(&Temp{}), ValidateField(true))
		if filter == nil {
			t.Error("Filter should not be nil")
		}

		if filter.kyte.source == nil {
			t.Error("Filter.kyte should not be nil")
		}

		if filter.kyte.checkField != true {
			t.Error("Filter.kyte.checkField should be true")
		}
	})

	t.Run("with validate field is false", func(t *testing.T) {
		filter := Filter(ValidateField(false))
		if filter == nil {
			t.Error("Filter should not be nil")
		}

		if filter.kyte.source != nil {
			t.Error("Filter.kyte should be nil")
		}

		if filter.kyte.checkField != false {
			t.Error("Filter.kyte.checkField should be false")
		}
	})

}

func TestFilter_Equal(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().Equal("name", "kyte").Build()
		if err != nil {
			t.Errorf("Filter.Equal should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Equal should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Equal should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$eq"] != "kyte" {
			t.Errorf("Filter.Equal should return value map[$eq:kyte], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).Equal(&temp.Name, "kyte").Build()
		if err != nil {
			t.Errorf("Filter.Equal should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Equal should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Equal should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$eq"] != "kyte" {
			t.Errorf("Filter.Equal should return value map[$eq:kyte], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
			Age     int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			Equal(&temp.Name, "Joe").
			Equal(&temp.Surname, "Doe").
			Equal(&temp.Age, 10).
			Build()

		if err != nil {
			t.Errorf("Filter.Equal should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Equal should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$eq"] != "Joe" {
					t.Errorf("Filter.Equal should return value map[$eq:Joe], got %v", v.Value)
				}
			}

			if v.Key == "surname" {
				if v.Value.(bson.M)["$eq"] != "Doe" {
					t.Errorf("Filter.Equal should return value map[$eq:Doe], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$eq"] != 10 {
					t.Errorf("Filter.Equal should return value map[$eq:10], got %v", v.Value)
				}
			}
		}

	})
}

func TestFilter_NotEqual(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().NotEqual("name", "kyte").Build()
		if err != nil {
			t.Errorf("Filter.NotEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NotEqual should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.NotEqual should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$ne"] != "kyte" {
			t.Errorf("Filter.NotEqual should return value map[$ne:kyte], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).NotEqual(&temp.Name, "kyte").Build()
		if err != nil {
			t.Errorf("Filter.NotEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NotEqual should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.NotEqual should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$ne"] != "kyte" {
			t.Errorf("Filter.NotEqual should return value map[$ne:kyte], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
			Age     int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			NotEqual(&temp.Name, "Joe").
			NotEqual(&temp.Surname, "Doe").
			NotEqual(&temp.Age, 10).
			Build()

		if err != nil {
			t.Errorf("Filter.NotEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NotEqual should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$ne"] != "Joe" {
					t.Errorf("Filter.NotEqual should return value map[$ne:Joe], got %v", v.Value)
				}
			}

			if v.Key == "surname" {
				if v.Value.(bson.M)["$ne"] != "Doe" {
					t.Errorf("Filter.NotEqual should return value map[$ne:Doe], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$ne"] != 10 {
					t.Errorf("Filter.NotEqual should return value map[$ne:10], got %v", v.Value)
				}
			}
		}
	})
}

func TestFilter_GreaterThan(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().GreaterThan("age", 10).Build()
		if err != nil {
			t.Errorf("Filter.GreaterThan should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.GreaterThan should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.GreaterThan should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$gt"] != 10 {
			t.Errorf("Filter.GreaterThan should return value map[$gt:10], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Age int `bson:"age"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).GreaterThan(&temp.Age, 10).Build()
		if err != nil {
			t.Errorf("Filter.GreaterThan should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.GreaterThan should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.GreaterThan should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$gt"] != 10 {
			t.Errorf("Filter.GreaterThan should return value map[$gt:10], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			GreaterThan(&temp.Name, "Joe").
			GreaterThan(&temp.Age, 10).
			Build()

		if err != nil {
			t.Errorf("Filter.GreaterThan should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.GreaterThan should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$gt"] != "Joe" {
					t.Errorf("Filter.GreaterThan should return value map[$gt:Joe], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$gt"] != 10 {
					t.Errorf("Filter.GreaterThan should return value map[$gt:10], got %v", v.Value)
				}
			}
		}
	})
}

func TestFilter_GreaterThanOrEqual(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().GreaterThanOrEqual("age", 10).Build()
		if err != nil {
			t.Errorf("Filter.GreaterThanOrEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.GreaterThanOrEqual should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.GreaterThanOrEqual should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$gte"] != 10 {
			t.Errorf("Filter.GreaterThanOrEqual should return value map[$gte:10], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Age int `bson:"age"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).GreaterThanOrEqual(&temp.Age, 10).Build()
		if err != nil {
			t.Errorf("Filter.GreaterThanOrEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.GreaterThanOrEqual should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.GreaterThanOrEqual should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$gte"] != 10 {
			t.Errorf("Filter.GreaterThanOrEqual should return value map[$gte:10], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			GreaterThanOrEqual(&temp.Name, "Joe").
			GreaterThanOrEqual(&temp.Age, 10).
			Build()

		if err != nil {
			t.Errorf("Filter.GreaterThanOrEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.GreaterThanOrEqual should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$gte"] != "Joe" {
					t.Errorf("Filter.GreaterThanOrEqual should return value map[$gte:Joe], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$gte"] != 10 {
					t.Errorf("Filter.GreaterThanOrEqual should return value map[$gte:10], got %v", v.Value)
				}
			}
		}
	})
}

func TestFilter_LessThan(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().LessThan("age", 10).Build()
		if err != nil {
			t.Errorf("Filter.LessThan should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.LessThan should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.LessThan should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$lt"] != 10 {
			t.Errorf("Filter.LessThan should return value map[$lt:10], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Age int `bson:"age"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).LessThan(&temp.Age, 10).Build()
		if err != nil {
			t.Errorf("Filter.LessThan should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.LessThan should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.LessThan should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$lt"] != 10 {
			t.Errorf("Filter.LessThan should return value map[$lt:10], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			LessThan(&temp.Name, "Joe").
			LessThan(&temp.Age, 10).
			Build()

		if err != nil {
			t.Errorf("Filter.LessThan should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.LessThan should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$lt"] != "Joe" {
					t.Errorf("Filter.LessThan should return value map[$lt:Joe], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$lt"] != 10 {
					t.Errorf("Filter.LessThan should return value map[$lt:10], got %v", v.Value)
				}
			}
		}
	})
}

func TestFilter_LessThanOrEqual(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().LessThanOrEqual("age", 10).Build()
		if err != nil {
			t.Errorf("Filter.LessThanOrEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.LessThanOrEqual should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.LessThanOrEqual should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$lte"] != 10 {
			t.Errorf("Filter.LessThanOrEqual should return value map[$lte:10], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Age int `bson:"age"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).LessThanOrEqual(&temp.Age, 10).Build()
		if err != nil {
			t.Errorf("Filter.LessThanOrEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.LessThanOrEqual should not return nil")
		}

		if q[0].Key != "age" {
			t.Errorf("Filter.LessThanOrEqual should return key age, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$lte"] != 10 {
			t.Errorf("Filter.LessThanOrEqual should return value map[$lte:10], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			LessThanOrEqual(&temp.Name, "Joe").
			LessThanOrEqual(&temp.Age, 10).
			Build()

		if err != nil {
			t.Errorf("Filter.LessThanOrEqual should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.LessThanOrEqual should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$lte"] != "Joe" {
					t.Errorf("Filter.LessThanOrEqual should return value map[$lte:Joe], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$lte"] != 10 {
					t.Errorf("Filter.LessThanOrEqual should return value map[$lte:10], got %v", v.Value)
				}
			}
		}
	})
}

func TestFilter_In(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		arr := []string{"kyte", "joe"}
		q, err := Filter().In("name", arr).Build()
		if err != nil {
			t.Errorf("Filter.In should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.In should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.In should return key name, got %v", q[0].Key)
		}

		if !reflect.DeepEqual(q[0].Value.(bson.M)["$in"], arr) {
			t.Errorf("Filter.In should return value %v, got %v", arr, q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name []string `bson:"name"`
		}
		var temp Temp
		arr := []string{"kyte", "joe"}
		q, err := Filter(Source(&temp)).In(&temp.Name, arr).Build()
		if err != nil {
			t.Errorf("Filter.In should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.In should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.In should return key name, got %v", q[0].Key)
		}

		if !reflect.DeepEqual(q[0].Value.(bson.M)["$in"], arr) {
			t.Errorf("Filter.In should return value %v, got %v", arr, q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name []string `bson:"name"`
			Age  []int    `bson:"age"`
		}

		var temp Temp
		arrName := []string{"kyte", "joe"}
		arrAge := []int{10, 20}
		q, err := Filter(Source(&temp)).
			In(&temp.Name, arrName).
			In(&temp.Age, arrAge).
			Build()

		if err != nil {
			t.Errorf("Filter.In should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.In should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if !reflect.DeepEqual(v.Value.(bson.M)["$in"], arrName) {
					t.Errorf("Filter.In should return value %v, got %v", arrName, v.Value)
				}
			}

			if v.Key == "age" {
				if !reflect.DeepEqual(v.Value.(bson.M)["$in"], arrAge) {
					t.Errorf("Filter.In should return value %v, got %v", arrAge, v.Value)
				}
			}
		}
	})

}

func TestFilter_NotIn(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		arr := []string{"kyte", "joe"}
		q, err := Filter().NotIn("name", arr).Build()
		if err != nil {
			t.Errorf("Filter.NotIn should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NotIn should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.NotIn should return key name, got %v", q[0].Key)
		}

		if !reflect.DeepEqual(q[0].Value.(bson.M)["$nin"], arr) {
			t.Errorf("Filter.NotIn should return value %v, got %v", arr, q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name []string `bson:"name"`
		}
		var temp Temp
		arr := []string{"kyte", "joe"}
		q, err := Filter(Source(&temp)).NotIn(&temp.Name, arr).Build()
		if err != nil {
			t.Errorf("Filter.NotIn should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NotIn should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.NotIn should return key name, got %v", q[0].Key)
		}

		if !reflect.DeepEqual(q[0].Value.(bson.M)["$nin"], arr) {
			t.Errorf("Filter.NotIn should return value %v, got %v", arr, q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name []string `bson:"name"`
			Age  []int    `bson:"age"`
		}

		var temp Temp
		arrName := []string{"kyte", "joe"}
		arrAge := []int{10, 20}
		q, err := Filter(Source(&temp)).
			NotIn(&temp.Name, arrName).
			NotIn(&temp.Age, arrAge).
			Build()

		if err != nil {
			t.Errorf("Filter.NotIn should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NotIn should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if !reflect.DeepEqual(v.Value.(bson.M)["$nin"], arrName) {
					t.Errorf("Filter.NotIn should return value %v, got %v", arrName, v.Value)
				}
			}

			if v.Key == "age" {
				if !reflect.DeepEqual(v.Value.(bson.M)["$nin"], arrAge) {
					t.Errorf("Filter.NotIn should return value %v, got %v", arrAge, v.Value)
				}
			}
		}
	})

}

func TestFilter_And(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().And(
			Filter().
				Equal("name", "kyte").
				Equal("surname", "joe"),
		).Build()

		if err != nil {
			t.Errorf("Filter.AND should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.AND should not return nil")
		}

		if q[0].Key != "$and" {
			t.Errorf("Filter.AND should return key $and, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != "kyte" {
			t.Errorf("Filter.AND should return value map[$eq:kyte], got %v", q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != "joe" {
			t.Errorf("Filter.AND should return value map[$eq:joe], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
		}
		var temp Temp
		name := "kyte"
		surname := "joe"
		q, err := Filter(Source(&temp)).And(
			Filter().
				Equal(&temp.Name, name).
				Equal(&temp.Surname, surname),
		).Build()

		if err != nil {
			t.Errorf("Filter.AND should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.AND should not return nil")
		}

		if q[0].Key != "$and" {
			t.Errorf("Filter.AND should return key $and, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != name {
			t.Errorf("Filter.AND should return value map[$eq:%v], got %v", name, q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != surname {
			t.Errorf("Filter.AND should return value map[$eq:%v], got %v", surname, q[0].Value)
		}

	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
			Age     int    `bson:"age"`
		}
		var temp Temp
		name := "kyte"
		surname := "joe"
		age := 10
		q, err := Filter(Source(&temp)).
			And(
				Filter().
					Equal(&temp.Name, name).
					Equal(&temp.Surname, surname),
			).
			And(
				Filter().
					GreaterThan(&temp.Age, age),
			).
			Build()

		if err != nil {
			t.Errorf("Filter.AND should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.AND should not return nil")
		}

		if q[0].Key != "$and" {
			t.Errorf("Filter.AND should return key $and, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != name {
			t.Errorf("Filter.AND should return value map[$eq:%v], got %v", name, q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != surname {
			t.Errorf("Filter.AND should return value map[$eq:%v], got %v", surname, q[0].Value)
		}

		if q[1].Key != "$and" {
			t.Errorf("Filter.AND should return key $and, got %v", q[1].Key)
		}

		if q[1].Value.(bson.A)[0].(bson.M)["age"].(bson.M)["$gt"] != age {
			t.Errorf("Filter.AND should return value map[$gt:%v], got %v", age, q[1].Value)
		}
	})

	t.Run("error on filter", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
		}
		var temp Temp
		name := "kyte"
		_, err := Filter(Source(&temp)).And(
			Filter().
				Equal(nil, name),
		).Build()

		if err == nil {
			t.Error("Filter.AND should return error")
		}
	})
}

func TestFilter_Or(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().Or(
			Filter().
				Equal("name", "kyte").
				Equal("surname", "joe"),
		).Build()

		if err != nil {
			t.Errorf("Filter.OR should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.OR should not return nil")
		}

		if q[0].Key != "$or" {
			t.Errorf("Filter.OR should return key $or, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != "kyte" {
			t.Errorf("Filter.OR should return value map[$eq:kyte], got %v", q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != "joe" {
			t.Errorf("Filter.OR should return value map[$eq:joe], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
		}
		var temp Temp
		name := "kyte"
		surname := "joe"
		q, err := Filter(Source(&temp)).Or(
			Filter().
				Equal(&temp.Name, name).
				Equal(&temp.Surname, surname),
		).Build()

		if err != nil {
			t.Errorf("Filter.OR should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.OR should not return nil")
		}

		if q[0].Key != "$or" {
			t.Errorf("Filter.OR should return key $or, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != name {
			t.Errorf("Filter.OR should return value map[$eq:%v], got %v", name, q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != surname {
			t.Errorf("Filter.OR should return value map[$eq:%v], got %v", surname, q[0].Value)
		}

	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
			Age     int    `bson:"age"`
		}
		var temp Temp
		name := "kyte"
		surname := "joe"
		age := 10
		q, err := Filter(Source(&temp)).
			Or(
				Filter().
					Equal(&temp.Name, name).
					Equal(&temp.Surname, surname),
			).
			Or(
				Filter().
					GreaterThan(&temp.Age, age),
			).
			Build()

		if err != nil {
			t.Errorf("Filter.OR should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.OR should not return nil")
		}

		if q[0].Key != "$or" {
			t.Errorf("Filter.OR should return key $or, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != name {
			t.Errorf("Filter.OR should return value map[$eq:%v], got %v", name, q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != surname {
			t.Errorf("Filter.OR should return value map[$eq:%v], got %v", surname, q[0].Value)
		}

		if q[1].Key != "$or" {
			t.Errorf("Filter.OR should return key $or, got %v", q[1].Key)
		}

		if q[1].Value.(bson.A)[0].(bson.M)["age"].(bson.M)["$gt"] != age {
			t.Errorf("Filter.OR should return value map[$gt:%v], got %v", age, q[1].Value)
		}
	})

	t.Run("error on filter", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
		}
		var temp Temp
		name := "kyte"
		_, err := Filter(Source(&temp)).Or(
			Filter().
				Equal(nil, name),
		).Build()

		if err == nil {
			t.Error("Filter.OR should return error")
		}
	})
}

func TestFilter_NOR(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().NOR(
			Filter().
				Equal("name", "kyte").
				Equal("surname", "joe"),
		).Build()

		if err != nil {
			t.Errorf("Filter.NOR should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NOR should not return nil")
		}

		if q[0].Key != "$nor" {
			t.Errorf("Filter.NOR should return key $nor, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != "kyte" {
			t.Errorf("Filter.NOR should return value map[$eq:kyte], got %v", q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != "joe" {
			t.Errorf("Filter.NOR should return value map[$eq:joe], got %v", q[0].Value)
		}
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
		}
		var temp Temp
		name := "kyte"
		surname := "joe"
		q, err := Filter(Source(&temp)).NOR(
			Filter().
				Equal(&temp.Name, name).
				Equal(&temp.Surname, surname),
		).Build()

		if err != nil {
			t.Errorf("Filter.NOR should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NOR should not return nil")
		}

		if q[0].Key != "$nor" {
			t.Errorf("Filter.NOR should return key $nor, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != name {
			t.Errorf("Filter.NOR should return value map[$eq:%v], got %v", name, q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != surname {
			t.Errorf("Filter.NOR should return value map[$eq:%v], got %v", surname, q[0].Value)
		}

	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
			Age     int    `bson:"age"`
		}
		var temp Temp
		name := "kyte"
		surname := "joe"
		age := 10
		q, err := Filter(Source(&temp)).
			NOR(
				Filter().
					Equal(&temp.Name, name).
					Equal(&temp.Surname, surname),
			).
			NOR(
				Filter().
					GreaterThan(&temp.Age, age),
			).
			Build()

		if err != nil {
			t.Errorf("Filter.NOR should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.NOR should not return nil")
		}

		if q[0].Key != "$nor" {
			t.Errorf("Filter.NOR should return key $nor, got %v", q[0].Key)
		}

		if q[0].Value.(bson.A)[0].(bson.M)["name"].(bson.M)["$eq"] != name {
			t.Errorf("Filter.NOR should return value map[$eq:%v], got %v", name, q[0].Value)
		}

		if q[0].Value.(bson.A)[1].(bson.M)["surname"].(bson.M)["$eq"] != surname {
			t.Errorf("Filter.NOR should return value map[$eq:%v], got %v", surname, q[0].Value)
		}

		if q[1].Key != "$nor" {
			t.Errorf("Filter.NOR should return key $nor, got %v", q[1].Key)
		}

		if q[1].Value.(bson.A)[0].(bson.M)["age"].(bson.M)["$gt"] != age {
			t.Errorf("Filter.NOR should return value map[$gt:%v], got %v", age, q[1].Value)
		}
	})

	t.Run("error on filter", func(t *testing.T) {
		type Temp struct {
			Name    string `bson:"name"`
			Surname string `bson:"surname"`
		}
		var temp Temp
		name := "kyte"
		_, err := Filter(Source(&temp)).NOR(
			Filter().
				Equal(nil, name),
		).Build()

		if err == nil {
			t.Error("Filter.NOR should return error")
		}
	})
}

func TestFilter_Regex(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().Regex("name", regexp.MustCompile("kyte")).Build()
		if err != nil {
			t.Errorf("Filter.Regex should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Regex should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Regex should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$regex"] != "kyte" {
			t.Errorf("Filter.Regex should return value map[$regex:kyte], got %v", q[0].Value)
		}
		fmt.Println(q[0].Value)
	})

	t.Run("with source", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).Regex(&temp.Name, regexp.MustCompile("kyte")).Build()
		if err != nil {
			t.Errorf("Filter.Regex should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Regex should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Regex should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$regex"] != "kyte" {
			t.Errorf("Filter.Regex should return value map[$regex:kyte], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			Regex(&temp.Name, regexp.MustCompile("kyte")).
			Regex(&temp.Age, regexp.MustCompile("10")).
			Build()

		if err != nil {
			t.Errorf("Filter.Regex should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Regex should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$regex"] != "kyte" {
					t.Errorf("Filter.Regex should return value map[$regex:kyte], got %v", v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$regex"] != "10" {
					t.Errorf("Filter.Regex should return value map[$regex:10], got %v", v.Value)
				}
			}
		}
	})

	t.Run("error on filter", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		_, err := Filter(Source(&temp)).Regex(nil, regexp.MustCompile("kyte")).Build()

		if err == nil {
			t.Error("Filter.Regex should return error")
		}
	})

	t.Run("error on regex", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		_, err := Filter(Source(&temp)).Regex(&temp.Name, nil).Build()

		if err != ErrRegexCannotBeNil {
			t.Errorf("Filter.Regex should return error %v, got %v", ErrRegexCannotBeNil, err)
		}
	})

	t.Run("with regex options", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).Regex(&temp.Name, regexp.MustCompile("kyte"), "i").Build()
		if err != nil {
			t.Errorf("Filter.Regex should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Regex should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Regex should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$regex"] != "kyte" {
			t.Errorf("Filter.Regex should return value map[$regex:kyte], got %v", q[0].Value)
		}

		if q[0].Value.(bson.M)["$options"] != "i" {
			t.Errorf("Filter.Regex should return value map[$options:i], got %v", q[0].Value)
		}
	})

	t.Run("with regex multiple options", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}
		var temp Temp
		q, err := Filter(Source(&temp)).Regex(&temp.Name, regexp.MustCompile("kyte"), "s", "i").Build()
		if err != nil {
			t.Errorf("Filter.Regex should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Regex should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Regex should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$regex"] != "kyte" {
			t.Errorf("Filter.Regex should return value map[$regex:kyte], got %v", q[0].Value)
		}

		if q[0].Value.(bson.M)["$options"] != "s" {
			t.Errorf("Filter.Regex should return value map[$options:s], got %v", q[0].Value)
		}
	})

}
