package kyte

import (
	"reflect"
	"regexp"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestFilter_Exists(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().Exists("name", true).Build()
		if err != nil {
			t.Errorf("Filter.Exists should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Exists should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Exists should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$exists"] != true {
			t.Errorf("Filter.Exists should return value map[$exists:true], got %v", q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		isNameExists := true
		isAgeExists := false
		q, err := Filter(Source(&temp)).
			Exists(&temp.Name, isNameExists).
			Exists(&temp.Age, isAgeExists).
			Build()

		if err != nil {
			t.Errorf("Filter.Exists should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Exists should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$exists"] != isNameExists {
					t.Errorf("Filter.Exists should return value %v, got %v", isNameExists, v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$exists"] != isAgeExists {
					t.Errorf("Filter.Exists should return value %v, got %v", isAgeExists, v.Value)
				}
			}
		}
	})
}

func TestFilter_Type(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		q, err := Filter().Type("name", bson.TypeString).Build()
		if err != nil {
			t.Errorf("Filter.Type should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Type should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Type should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$type"].([]bsontype.Type)[0] != bson.TypeString {
			t.Errorf("Filter.Type should return value  %v, got %v", bson.TypeString, q[0].Value.(bson.M)["$type"].([]bsontype.Type)[0])
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		q, err := Filter(Source(&temp)).
			Type(&temp.Name, bson.TypeString).
			Type(&temp.Age, bson.TypeInt32).
			Build()

		if err != nil {
			t.Errorf("Filter.Type should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Type should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$type"].([]bsontype.Type)[0] != bson.TypeString {
					t.Errorf("Filter.Type should return value %v, got %v", bson.TypeString, v.Value.(bson.M)["$type"].([]bsontype.Type)[0])
				}
			}
			if v.Key == "age" {
				if v.Value.(bson.M)["$type"].([]bsontype.Type)[0] != bson.TypeInt32 {
					t.Errorf("Filter.Type should return value  %v, got %v", bson.TypeInt32, v.Value.(bson.M)["$type"].([]bsontype.Type)[0])
				}
			}
		}
	})

	t.Run("mutliple types", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}

		var temp Temp
		type1 := bson.TypeString
		type2 := bson.TypeInt32
		q, err := Filter(Source(&temp)).Type(&temp.Name, type1, type2).Build()
		if err != nil {
			t.Errorf("Filter.Type should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Type should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Type should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$type"].([]bsontype.Type)[0] != type1 {
			t.Errorf("Filter.Type should return value  %v, got %v", type1, q[0].Value.(bson.M)["$type"].([]bsontype.Type)[0])
		}

		if q[0].Value.(bson.M)["$type"].([]bsontype.Type)[1] != type2 {
			t.Errorf("Filter.Type should return value  %v, got %v", type2, q[0].Value.(bson.M)["$type"].([]bsontype.Type)[1])
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()

		t.Run("zero type", func(t *testing.T) {
			type Temp struct {
				Name string `bson:"name"`
			}

			var temp Temp
			_, err := Filter(Source(&temp)).Type(&temp.Name).Build()

			if err != ErrInvalidBsonType {
				t.Errorf("Filter.Type should return error %v, got %v", ErrInvalidBsonType, err)
			}
		})

		t.Run("invalid type", func(t *testing.T) {
			type Temp struct {
				Name string `bson:"name"`
			}

			var temp Temp
			_, err := Filter(Source(&temp)).Type(&temp.Name, 100).Build()

			if err != ErrInvalidBsonType {
				t.Errorf("Filter.Type should return error %v, got %v", ErrInvalidBsonType, err)
			}
		})
	})
}

func TestFilter_Mod(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		divisor := 10
		remainder := 1
		q, err := Filter().Mod("name", divisor, remainder).Build()
		if err != nil {
			t.Errorf("Filter.Mod should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Mod should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Mod should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$mod"].(primitive.A)[0].(int) != divisor {
			t.Errorf("Filter.Mod should return value map[$mod:[%v %v]], got %v", divisor, remainder, q[0].Value)
		}

		if q[0].Value.(bson.M)["$mod"].(primitive.A)[1].(int) != remainder {
			t.Errorf("Filter.Mod should return value map[$mod:[%v %v]], got %v", divisor, remainder, q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		divisor := 10
		divisor1 := 20
		remainder := 1
		remainder1 := 2
		q, err := Filter(Source(&temp)).
			Mod(&temp.Name, divisor, remainder).
			Mod(&temp.Age, divisor1, remainder1).
			Build()

		if err != nil {
			t.Errorf("Filter.Mod should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Mod should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$mod"].(primitive.A)[0].(int) != divisor {
					t.Errorf("Filter.Mod should return value map[$mod:[%v %v]], got %v", divisor, remainder, v.Value)
				}

				if v.Value.(bson.M)["$mod"].(primitive.A)[1].(int) != remainder {
					t.Errorf("Filter.Mod should return value map[$mod:[%v %v]], got %v", divisor, remainder, v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$mod"].(primitive.A)[0].(int) != divisor1 {
					t.Errorf("Filter.Mod should return value map[$mod:[%v %v]], got %v", divisor1, remainder1, v.Value)
				}

				if v.Value.(bson.M)["$mod"].(primitive.A)[1].(int) != remainder1 {
					t.Errorf("Filter.Mod should return value map[$mod:[%v %v]], got %v", divisor1, remainder1, v.Value)
				}
			}
		}
	})
}

func TestFilter_Where(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		fn := "this.name == 'kyte'"
		q, err := Filter().Where(fn).Build()
		if err != nil {
			t.Errorf("Filter.Where should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Where should not return nil")
		}

		if q[0].Key != "$where" {
			t.Errorf("Filter.Where should return key $where, got %v", q[0].Key)
		}

		if q[0].Value.(string) != fn {
			t.Errorf("Filter.Where should return value %v, got %v", fn, q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
			Age  int    `bson:"age"`
		}

		var temp Temp
		fn := "this.name == 'kyte'"
		fn1 := "this.age == 10"
		q, err := Filter(Source(&temp)).
			Where(fn).
			Where(fn1).
			Build()

		if err != nil {
			t.Errorf("Filter.Where should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Where should not return nil")
		}

		fns := []string{fn, fn1}

		for _, v := range q {
			if !contains(fns, v.Value.(string)) {
				t.Errorf("Filter.Where should return value %v, got %v", fns, v.Value)
			}
		}
	})
}

func TestFilter_All(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		arr := []string{"kyte", "joe"}
		q, err := Filter().All("name", arr).Build()
		if err != nil {
			t.Errorf("Filter.All should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.All should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.All should return key name, got %v", q[0].Key)
		}

		if !reflect.DeepEqual(q[0].Value.(bson.M)["$all"], arr) {
			t.Errorf("Filter.All should return value %v, got %v", arr, q[0].Value)
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
			All(&temp.Name, arrName).
			All(&temp.Age, arrAge).
			Build()

		if err != nil {
			t.Errorf("Filter.All should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.All should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if !reflect.DeepEqual(v.Value.(bson.M)["$all"], arrName) {
					t.Errorf("Filter.All should return value %v, got %v", arrName, v.Value)
				}
			}

			if v.Key == "age" {
				if !reflect.DeepEqual(v.Value.(bson.M)["$all"], arrAge) {
					t.Errorf("Filter.All should return value %v, got %v", arrAge, v.Value)
				}
			}
		}

	})

	t.Run("not slice", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}

		var temp Temp
		_, err := Filter(Source(&temp)).All(&temp.Name, "arr").Build()

		if err != ErrValueMustBeSlice {
			t.Errorf("Filter.All should return error %v, got %v", ErrValueMustBeSlice, err)
		}
	})
}

func TestFilter_Size(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {
		size := 10
		q, err := Filter().Size("name", size).Build()
		if err != nil {
			t.Errorf("Filter.Size should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Size should not return nil")
		}

		if q[0].Key != "name" {
			t.Errorf("Filter.Size should return key name, got %v", q[0].Key)
		}

		if q[0].Value.(bson.M)["$size"] != size {
			t.Errorf("Filter.Size should return value %v, got %v", size, q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name []string `bson:"name"`
			Age  []int    `bson:"age"`
		}

		var temp Temp
		sizeName := 10
		sizeAge := 20
		q, err := Filter(Source(&temp)).
			Size(&temp.Name, sizeName).
			Size(&temp.Age, sizeAge).
			Build()

		if err != nil {
			t.Errorf("Filter.Size should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.Size should not return nil")
		}

		for _, v := range q {
			if v.Key == "name" {
				if v.Value.(bson.M)["$size"] != sizeName {
					t.Errorf("Filter.Size should return value %v, got %v", sizeName, v.Value)
				}
			}

			if v.Key == "age" {
				if v.Value.(bson.M)["$size"] != sizeAge {
					t.Errorf("Filter.Size should return value %v, got %v", sizeAge, v.Value)
				}
			}
		}
	})
}

func TestFilter_JSONSchema(t *testing.T) {
	t.Parallel()

	t.Run("without source", func(t *testing.T) {

		schema := bson.M{
			"type": "object",
			"properties": bson.M{
				"name": bson.M{
					"type": "string",
				},
			},
		}

		q, err := Filter().JSONSchema(schema).Build()
		if err != nil {
			t.Errorf("Filter.JSONSchema should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.JSONSchema should not return nil")
		}

		if q[0].Key != "$jsonSchema" {
			t.Errorf("Filter.JSONSchema should return key $jsonSchema, got %v", q[0].Key)
		}

		if !reflect.DeepEqual(q[0].Value, schema) {
			t.Errorf("Filter.JSONSchema should return value %v, got %v", schema, q[0].Value)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}

		var temp Temp
		schema1 := bson.M{
			"type": "object",
			"properties": bson.M{
				"name": bson.M{
					"type": "string",
				},
			},
		}

		schema2 := bson.M{
			"type": "object",
			"properties": bson.M{
				"age": bson.M{
					"type": "int",
				},
			},
		}

		q, err := Filter(Source(&temp)).JSONSchema(schema1).JSONSchema(schema2).Build()
		if err != nil {
			t.Errorf("Filter.JSONSchema should not return error: %v", err)
		}

		if q == nil {
			t.Error("Filter.JSONSchema should not return nil")
		}

		// schemas := []bson.M{schema1, schema2}
		for _, v := range q {
			if v.Value.(bson.M)["properties"].(bson.M)["name"] != nil {
				if !reflect.DeepEqual(v.Value.(bson.M), schema1) {
					t.Errorf("Filter.JSONSchema should return value %v, got %v", schema1, v.Value)
				}
			}

			if v.Value.(bson.M)["properties"].(bson.M)["age"] != nil {
				if !reflect.DeepEqual(v.Value.(bson.M), schema2) {
					t.Errorf("Filter.JSONSchema should return value %v, got %v", schema2, v.Value)
				}
			}
		}
	})
}
