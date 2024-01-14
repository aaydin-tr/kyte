package main

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	eq  = "$eq"
	ne  = "$ne"
	gt  = "$gt"
	gte = "$gte"
	lt  = "$lt"
	lte = "$lte"
	in  = "$in"
	nin = "$nin"

	and = "$and"
	or  = "$or"
	// TODO: implement
	// not = "$not"
	// nor = "$nor"
)

type filter struct {
	Kyte
	query bson.M
}

/*
Filter creates a new filter instance.
Source must be a pointer of a struct.
CheckField is true by default, it will check if the field is valid for query based on the source struct bson tags.
*/
func Filter(source any, checkField ...bool) *filter {
	kyte := newKyte(source)

	if len(checkField) > 0 {
		kyte.fieldCheck = checkField[0]
	}

	return &filter{
		Kyte:  *kyte,
		query: bson.M{},
	}
}

/*
Equal use mongo [$eq] operator to compare field and value.

	Filter(source).
		Equal("name", "John") // {"name": {"$eq": "John"}}

[$eq]: https://www.mongodb.com/docs/manual/reference/operator/query/eq/#mongodb-query-op.-eq
*/
func (f *filter) Equal(field any, value any) *filter {
	return f.set(eq, field, value)
}

/*
NotEqual use mongo [$ne] operator to compare field and value.

	Filter(source).
		NotEqual("name", "John") // {"name": {"$ne": "John"}}

[$ne]: https://www.mongodb.com/docs/manual/reference/operator/query/ne/#mongodb-query-op.-ne
*/
func (f *filter) NotEqual(field any, value any) *filter {
	return f.set(ne, field, value)
}

/*
GreaterThan use mongo [$gt] operator to compare field and value.

	Filter(source).
		GreaterThan("age", 18) // {"age": {"$gt": 18}}

[$gt]: https://www.mongodb.com/docs/manual/reference/operator/query/gt/#mongodb-query-op.-gt
*/
func (f *filter) GreaterThan(field any, value any) *filter {
	return f.set(gt, field, value)
}

/*
GreaterThanOrEqual use mongo [$gte] operator to compare field and value.

Example: GreaterThanOrEqual("age", 18) => {"age": {"$gte": 18}}

[$gte]: https://www.mongodb.com/docs/manual/reference/operator/query/gte/#mongodb-query-op.-gte
*/
func (f *filter) GreaterThanOrEqual(field any, value any) *filter {
	return f.set(gte, field, value)
}

/*
LessThan use mongo [$lt] operator to compare field and value.

	Filter(source).
		LessThan("age", 18) // {"age": {"$lt": 18}}

[$lt]: https://www.mongodb.com/docs/manual/reference/operator/query/lt/#mongodb-query-op.-lt
*/
func (f *filter) LessThan(field any, value any) *filter {
	return f.set(lt, field, value)
}

/*
LessThanOrEqual use mongo [$lte] operator to compare field and value.

	Filter(source).
		LessThanOrEqual("age", 18) // {"age": {"$lte": 18}}

[$lte]: https://www.mongodb.com/docs/manual/reference/operator/query/lte/#mongodb-query-op.-lte
*/
func (f *filter) LessThanOrEqual(field any, value any) *filter {
	return f.set(lte, field, value)
}

/*
In use mongo [$in] operator to compare field and value.

	Filter(source).
		In("name", []string{"John", "Jane"}) // {"name": {"$in": ["John", "Jane"]}}

[$in]: https://www.mongodb.com/docs/manual/reference/operator/query/in/#mongodb-query-op.-in
*/
func (f *filter) In(field any, value any) *filter {
	return f.set(in, field, value)
}

/*
NotIn use mongo [$nin] operator to compare field and value.

	Filter(source).
		NotIn("name", []string{"John", "Jane"}) // {"name": {"$nin": ["John", "Jane"]}}

[$nin]: https://www.mongodb.com/docs/manual/reference/operator/query/nin/#mongodb-query-op.-nin
*/
func (f *filter) NotIn(field any, value any) *filter {
	return f.set(nin, field, value)
}

/*
And use mongo [$and] logical query operator to combine multiple query expressions.

	Filter(source).
		Equal("name", "John").
		And(
			Filter(source).
				Equal("age", 18).
				Equal("surname", "Doe"),
		)

[$and]: https://www.mongodb.com/docs/manual/reference/operator/query/and/#mongodb-query-op.-and
*/
func (f *filter) And(filter *filter) *filter {
	query, err := filter.Build()
	if err != nil {
		f.setError(err)
		return f
	}

	if f.query[and] == nil {
		f.query[and] = bson.A{query}
	} else {
		f.query[and] = append(f.query[and].(bson.A), query)
	}

	return f
}

/*
Or use mongo [$or] logical query operator to combine multiple query expressions.

	Filter(source).
		Equal("name", "John").
		Or(
			Filter(source).
				Equal("age", 18).
				Equal("surname", "Doe"),
		)

[$or]: https://www.mongodb.com/docs/manual/reference/operator/query/or/#mongodb-query-op.-or
*/
func (f *filter) Or(filter *filter) *filter {
	query, err := filter.Build()
	if err != nil {
		f.setError(err)
		return f
	}

	if f.query[or] == nil {
		f.query[or] = bson.A{query}
	} else {
		f.query[or] = append(f.query[or].(bson.A), query)
	}

	return f
}

func (f *filter) set(operator string, field any, value any) *filter {
	fieldName, err := f.validateQueryFieldAndValue(field, value)
	if err != nil {
		f.setError(err)
		return f
	}

	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	if operator == in || operator == nin {
		if valueType.Kind() != reflect.Slice {
			value = bson.A{value}
		}

		f.query[fieldName] = bson.M{operator: value}
	} else {
		f.query[fieldName] = bson.M{operator: value}
	}
	return f
}

/*
Build returns the query as bson.M. If there is an error, it will return nil and the first error.
*/
func (f *filter) Build() (bson.M, error) {
	if f.hasErrors() {
		return nil, f.errs[0]
	}

	return f.query, nil
}
