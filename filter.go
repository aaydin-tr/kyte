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

type operation struct {
	operator string
	field    any
	value    any
}

type filter struct {
	kyte       Kyte
	query      bson.D
	operations []operation
}

type FilterOptions struct {
	// Source is the struct that will be used to check if the field is valid for query based on the struct bson tags.
	Source any

	// CheckField is true by default, it will check if the field is valid for query based on the source struct bson tags.
	CheckField bool
}

type FilterOption func(*FilterOptions)

func WithCheckField(checkField bool) FilterOption {
	return func(o *FilterOptions) {
		o.CheckField = checkField
	}
}

func WithSource(source any) FilterOption {
	return func(o *FilterOptions) {
		o.Source = source
	}
}

/*
Filter creates a new filter instance.
*/
func Filter(opts ...FilterOption) *filter {
	options := &FilterOptions{}
	for _, opt := range opts {
		opt(options)
	}

	kyte := newKyte(options.Source, options.CheckField)

	return &filter{
		kyte:  *kyte,
		query: bson.D{},
	}
}

/*
Equal use mongo [$eq] operator to compare field and value.

	Filter().
		Equal("name", "John") // {"name": {"$eq": "John"}}

[$eq]: https://www.mongodb.com/docs/manual/reference/operator/query/eq/#mongodb-query-op.-eq
*/
func (f *filter) Equal(field any, value any) *filter {
	return f.set(eq, field, value)
}

/*
NotEqual use mongo [$ne] operator to compare field and value.

	Filter().
		NotEqual("name", "John") // {"name": {"$ne": "John"}}

[$ne]: https://www.mongodb.com/docs/manual/reference/operator/query/ne/#mongodb-query-op.-ne
*/
func (f *filter) NotEqual(field any, value any) *filter {
	return f.set(ne, field, value)
}

/*
GreaterThan use mongo [$gt] operator to compare field and value.

	Filter().
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

	Filter().
		LessThan("age", 18) // {"age": {"$lt": 18}}

[$lt]: https://www.mongodb.com/docs/manual/reference/operator/query/lt/#mongodb-query-op.-lt
*/
func (f *filter) LessThan(field any, value any) *filter {
	return f.set(lt, field, value)
}

/*
LessThanOrEqual use mongo [$lte] operator to compare field and value.

	Filter().
		LessThanOrEqual("age", 18) // {"age": {"$lte": 18}}

[$lte]: https://www.mongodb.com/docs/manual/reference/operator/query/lte/#mongodb-query-op.-lte
*/
func (f *filter) LessThanOrEqual(field any, value any) *filter {
	return f.set(lte, field, value)
}

/*
In use mongo [$in] operator to compare field and value.

	Filter().
		In("name", []string{"John", "Jane"}) // {"name": {"$in": ["John", "Jane"]}}

[$in]: https://www.mongodb.com/docs/manual/reference/operator/query/in/#mongodb-query-op.-in
*/
func (f *filter) In(field any, value any) *filter {
	return f.set(in, field, value)
}

/*
NotIn use mongo [$nin] operator to compare field and value.

	Filter().
		NotIn("name", []string{"John", "Jane"}) // {"name": {"$nin": ["John", "Jane"]}}

[$nin]: https://www.mongodb.com/docs/manual/reference/operator/query/nin/#mongodb-query-op.-nin
*/
func (f *filter) NotIn(field any, value any) *filter {
	return f.set(nin, field, value)
}

/*
And use mongo [$and] logical query operator to combine multiple query expressions.

	Filter().
		Equal("name", "John").
		And(
			Filter().
				Equal("age", 18).
				Equal("surname", "Doe"),
		)

[$and]: https://www.mongodb.com/docs/manual/reference/operator/query/and/#mongodb-query-op.-and
*/
func (f *filter) And(filter *filter) *filter {
	if f.kyte.source != nil {
		filter.kyte.checkField = f.kyte.checkField
		filter.kyte.setSourceAndPrepareFields(f.kyte.source)
	}

	query, err := filter.Build()
	if err != nil {
		f.kyte.setError(err)
		return f
	}

	andQuery := bson.A{}
	for _, q := range query {
		andQuery = append(andQuery, bson.M{q.Key: q.Value})
	}

	f.query = append(f.query, bson.E{Key: and, Value: andQuery})
	return f
}

/*
Or use mongo [$or] logical query operator to combine multiple query expressions.

	Filter(source).
		Equal("name", "John").
		Or(
			Filter().
				Equal("age", 18).
				Equal("surname", "Doe"),
		)

[$or]: https://www.mongodb.com/docs/manual/reference/operator/query/or/#mongodb-query-op.-or
*/
func (f *filter) Or(filter *filter) *filter {
	if f.kyte.source != nil {
		filter.kyte.checkField = f.kyte.checkField
		filter.kyte.setSourceAndPrepareFields(f.kyte.source)
	}

	query, err := filter.Build()
	if err != nil {
		f.kyte.setError(err)
		return f
	}

	orQuery := bson.A{}
	for _, q := range query {
		orQuery = append(orQuery, bson.M{q.Key: q.Value})
	}

	f.query = append(f.query, bson.E{Key: or, Value: orQuery})
	return f
}

func (f *filter) set(operator string, field any, value any) *filter {
	if f.kyte.hasErrors() {
		return f
	}
	f.operations = append(f.operations, operation{operator: operator, field: field, value: value})
	return f
}

/*
Build returns the query as bson.M. If there is an error, it will return nil and the first error.
*/
func (f *filter) Build() (bson.D, error) {
	for _, opt := range f.operations {
		fieldName, err := f.kyte.validateQueryFieldAndValue(opt.field, opt.value)
		if err != nil {
			f.kyte.setError(err)
			break
		}

		valueType := reflect.TypeOf(opt.value)
		if valueType.Kind() == reflect.Ptr {
			opt.value = reflect.ValueOf(opt.value).Elem().Interface()
		}

		if opt.operator == in || opt.operator == nin {
			if valueType.Kind() != reflect.Slice {
				opt.value = bson.A{opt.value}
			}

		}
		f.query = append(f.query, bson.E{Key: fieldName, Value: bson.M{opt.operator: opt.value}})
	}

	if f.kyte.hasErrors() {
		return nil, f.kyte.errs[0]
	}

	return f.query, nil
}
