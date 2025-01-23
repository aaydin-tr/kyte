package kyte

import (
	"reflect"
	"regexp"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var (
	globalFilters []*filter
	globalMutex   sync.RWMutex
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
	nor = "$nor"

	regx        = "$regex"
	regxOptions = "$options"

	exists     = "$exists"
	_type      = "$type"
	mod        = "$mod"
	where      = "$where"
	all        = "$all"
	size       = "$size"
	jsonSchema = "$jsonSchema"

	// TODO implement Day 1
	// $elemMatch
	// $not

	// TODO: implement Later
	// $text
	// $expr
	// $geoIntersects
	// $geoWithin
	// $near
	// $nearSphere
	// $bitsAllClear
	// $bitsAllSet
	// $bitsAnyClear
	// $bitsAnySet
)

/*
AddGlobalFilter adds a filter that will be applied to all new filter instances.
This is useful for scenarios like multi-tenancy where certain conditions should
always be applied. This function is thread-safe.

Example:

	// Set up a global tenant filter
	kyte.AddGlobalFilter(kyte.Filter().Equal("tenantId", "123"))
*/
func AddGlobalFilter(f *filter) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalFilters = append(globalFilters, f)
}

/*
ClearGlobalFilters removes all registered global filters.
This function is thread-safe.
*/
func ClearGlobalFilters() {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalFilters = nil
}

/*
GetGlobalFilters returns a copy of the current global filters.
This function is thread-safe.
*/
func GetGlobalFilters() []*filter {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	// Return a copy to prevent external modifications
	result := make([]*filter, len(globalFilters))
	copy(result, globalFilters)
	return result
}

type filter struct {
	kyte       *kyte
	query      bson.D
	operations []operation

	isBuild bool
}

/*
Filter creates a new filter instance.
*/
func Filter(opts ...OptionFunc) *filter {
	options := &Options{validateField: true}
	for _, opt := range opts {
		opt(options)
	}

	kyte := newKyte(options.source, options.validateField)

	f := &filter{
		kyte:  kyte,
		query: bson.D{},
	}

	globalMutex.RLock()
	for _, globalFilter := range globalFilters {
		if globalQuery, err := globalFilter.Build(); err == nil {
			f.query = append(f.query, globalQuery...)
		}
	}
	globalMutex.RUnlock()

	return f
}

/*
Equal use mongo [$eq] operator to compare field and value.

	Filter().
		Equal("name", "John") // {"name": {"$eq": "John"}}

[$eq]: https://www.mongodb.com/docs/manual/reference/operator/query/eq/#mongodb-query-op.-eq
*/
func (f *filter) Equal(field any, value any) *filter {
	return f.set(eq, field, value, true)
}

/*
NotEqual use mongo [$ne] operator to compare field and value.

	Filter().
		NotEqual("name", "John") // {"name": {"$ne": "John"}}

[$ne]: https://www.mongodb.com/docs/manual/reference/operator/query/ne/#mongodb-query-op.-ne
*/
func (f *filter) NotEqual(field any, value any) *filter {
	return f.set(ne, field, value, true)
}

/*
GreaterThan use mongo [$gt] operator to compare field and value.

	Filter().
		GreaterThan("age", 18) // {"age": {"$gt": 18}}

[$gt]: https://www.mongodb.com/docs/manual/reference/operator/query/gt/#mongodb-query-op.-gt
*/
func (f *filter) GreaterThan(field any, value any) *filter {
	return f.set(gt, field, value, true)
}

/*
GreaterThanOrEqual use mongo [$gte] operator to compare field and value.

Example: GreaterThanOrEqual("age", 18) => {"age": {"$gte": 18}}

[$gte]: https://www.mongodb.com/docs/manual/reference/operator/query/gte/#mongodb-query-op.-gte
*/
func (f *filter) GreaterThanOrEqual(field any, value any) *filter {
	return f.set(gte, field, value, true)
}

/*
LessThan use mongo [$lt] operator to compare field and value.

	Filter().
		LessThan("age", 18) // {"age": {"$lt": 18}}

[$lt]: https://www.mongodb.com/docs/manual/reference/operator/query/lt/#mongodb-query-op.-lt
*/
func (f *filter) LessThan(field any, value any) *filter {
	return f.set(lt, field, value, true)
}

/*
LessThanOrEqual use mongo [$lte] operator to compare field and value.

	Filter().
		LessThanOrEqual("age", 18) // {"age": {"$lte": 18}}

[$lte]: https://www.mongodb.com/docs/manual/reference/operator/query/lte/#mongodb-query-op.-lte
*/
func (f *filter) LessThanOrEqual(field any, value any) *filter {
	return f.set(lte, field, value, true)
}

/*
In use mongo [$in] operator to compare field and value.

	Filter().
		In("name", []string{"John", "Jane"}) // {"name": {"$in": ["John", "Jane"]}}

[$in]: https://www.mongodb.com/docs/manual/reference/operator/query/in/#mongodb-query-op.-in
*/
func (f *filter) In(field any, value any) *filter {
	return f.set(in, field, value, true)
}

/*
NotIn use mongo [$nin] operator to compare field and value.

	Filter().
		NotIn("name", []string{"John", "Jane"}) // {"name": {"$nin": ["John", "Jane"]}}

[$nin]: https://www.mongodb.com/docs/manual/reference/operator/query/nin/#mongodb-query-op.-nin
*/
func (f *filter) NotIn(field any, value any) *filter {
	return f.set(nin, field, value, true)
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

	Filter().
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

/*
NOR use mongo [$nor] logical query operator to combine multiple query expressions.

	Filter(source).
		Equal("name", "John").
		NOR(
			Filter().
				Equal("age", 18).
				Equal("surname", "Doe"),
		)

[$nor]: https://www.mongodb.com/docs/manual/reference/operator/query/nor/#mongodb-query-op.-nor
*/
func (f *filter) NOR(filter *filter) *filter {
	if f.kyte.source != nil {
		filter.kyte.checkField = f.kyte.checkField
		filter.kyte.setSourceAndPrepareFields(f.kyte.source)
	}

	query, err := filter.Build()
	if err != nil {
		f.kyte.setError(err)
		return f
	}

	norQuery := bson.A{}
	for _, q := range query {
		norQuery = append(norQuery, bson.M{q.Key: q.Value})
	}

	f.query = append(f.query, bson.E{Key: nor, Value: norQuery})
	return f
}

/*
Regex use mongo [$regex] operator to compare field and value.

	Filter().
		Regex("name", regexp.MustCompile("^J")) // {"name": {"$regex": "^J"}}

	Filter().
		Regex("name", regexp.MustCompile("^J"), "i") // {"name": {"$regex": "^J", "$options": "i"}}

	Filter().
		Regex("name", regexp.MustCompile("^J"), "im") // {"name": {"$regex": "^J", "$options": "im"}}

[$regex]: https://www.mongodb.com/docs/manual/reference/operator/query/regex/#mongodb-query-op.-regex
*/
func (f *filter) Regex(field any, regex *regexp.Regexp, options ...string) *filter {
	if regex == nil {
		f.kyte.setError(ErrRegexCannotBeNil)
		return f
	}

	if len(options) == 0 {
		return f.set(regx, field, bson.M{regx: regex.String()}, true)
	}

	return f.set(regx, field, bson.M{regx: regex.String(), regxOptions: options[0]}, true)
}

/*
Exists use mongo [$exists] operator to check if the field exists.

	Filter().
		Exists("name", true) // {"name": {"$exists": true}}

[$exists]: https://www.mongodb.com/docs/manual/reference/operator/query/exists/#mongodb-query-op.-exists
*/
func (f *filter) Exists(field any, value bool) *filter {
	return f.set(exists, field, value, true)
}

/*
Type use mongo [$type] operator to check if the field is of the specified type. It accepts multiple types.

	Filter().
		Type("name", bsontype.String) // {"name": {"$type": "string"}}

	Filter().
		Type("name", bsontype.String, bsontype.Null) // {"name": {"$type": ["string", "null"]}}

[$type]: https://www.mongodb.com/docs/manual/reference/operator/query/type/#mongodb-query-op.-type
*/
func (f *filter) Type(field any, values ...bsontype.Type) *filter {
	if len(values) == 0 {
		f.kyte.setError(ErrInvalidBsonType)
		return f
	}

	for _, v := range values {
		if !v.IsValid() {
			f.kyte.setError(ErrInvalidBsonType)
			return f
		}
	}

	return f.set(_type, field, values, true)
}

/*
Mod use mongo [$mod] operator to check if the field is a multiple of a specified divisor.

	Filter().
		Mod("age", 2, 0) // {"age": {"$mod": [2, 0]}}

[$mod]: https://www.mongodb.com/docs/manual/reference/operator/query/mod/#mongodb-query-op.-mod
*/
func (f *filter) Mod(field any, divisor int, remainder int) *filter {
	return f.set(mod, field, bson.A{divisor, remainder}, true)
}

/*
Where use mongo [$where] operator to pass a javascript expression to the query system.

	Filter().
		Where("this.name === 'John'") // {"$where": "this.name === 'John'"}

[$where]: https://www.mongodb.com/docs/manual/reference/operator/query/where/#mongodb-query-op.-where
*/
func (f *filter) Where(js string) *filter {
	return f.set(where, nil, js, false)
}

/*
All use mongo [$all] operator to check if the field contains all of the specified elements.

	Filter().
		All("name", []string{"John", "Jane"}) // {"name": {"$all": ["John", "Jane"]}}

[$all]: https://www.mongodb.com/docs/manual/reference/operator/query/all/#mongodb-query-op.-all
*/
func (f *filter) All(field any, value any) *filter {
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		f.kyte.setError(ErrValueMustBeSlice)
		return f
	}

	return f.set(all, field, value, true)
}

/*
Size use mongo [$size] operator to check if the field is an array that contains a specific number of elements.

	Filter().
		Size("name", 2) // {"name": {"$size": 2}}

[$size]: https://www.mongodb.com/docs/manual/reference/operator/query/size/#mongodb-query-op.-size
*/
func (f *filter) Size(field any, value int) *filter {
	return f.set(size, field, value, true)
}

/*
JSONSchema use mongo [$jsonSchema] operator to validate documents against the given JSON Schema.

	Filter().
		JSONSchema(bson.M{"required": []string{"name"}}) // {"$jsonSchema": {"required": ["name"]}}

[$jsonSchema]: https://www.mongodb.com/docs/manual/reference/operator/query/jsonSchema/#mongodb-query-op.-jsonSchema
*/
func (f *filter) JSONSchema(schema bson.M) *filter {
	return f.set(jsonSchema, nil, schema, false)
}

/*
Raw use raw bson.D and directly append it to the query. It is useful for using operators that are not implemented in this package.
Raw will not provide any validation, so it is recommended to use it carefully.

	Filter().
		Raw(bson.D{{"name", "John"}}) // {"name": "John"}
*/
func (f *filter) Raw(query bson.D) *filter {
	f.query = append(f.query, query...)
	return f
}

/*
Build returns the query as bson.M. If there is an error, it will return nil and the first error.
*/
func (f *filter) Build() (bson.D, error) {
	if f.kyte.hasErrors() {
		return nil, f.kyte.err
	}

	if f.isBuild {
		return f.query, nil
	}

	for _, opt := range f.operations {
		err := f.kyte.validate(&opt)
		if err != nil {
			f.kyte.setError(err)
			break
		}

		if opt.operator == where {
			f.query = append(f.query, bson.E{Key: where, Value: opt.value})
			continue
		}

		if opt.operator == jsonSchema {
			f.query = append(f.query, bson.E{Key: jsonSchema, Value: opt.value})
			continue
		}

		fieldName, err := f.kyte.getFieldName(opt.field)
		if err != nil {
			f.kyte.setError(err)
			break
		}

		if opt.value != nil && reflect.TypeOf(opt.value).Kind() == reflect.Ptr {
			opt.value = reflect.ValueOf(opt.value).Elem().Interface()
		}

		if opt.operator == in || opt.operator == nin || opt.operator == _type {
			if opt.value != nil && reflect.TypeOf(opt.value).Kind() != reflect.Slice {
				opt.value = bson.A{opt.value}
			}
		}

		if opt.operator == regx {
			f.query = append(f.query, bson.E{Key: fieldName, Value: opt.value})
			continue
		}

		f.query = append(f.query, bson.E{Key: fieldName, Value: bson.M{opt.operator: opt.value}})
	}

	if f.kyte.hasErrors() {
		return nil, f.kyte.err
	}

	f.isBuild = true
	return f.query, nil
}

func (f *filter) set(operator string, field any, value any, isFieldRequired bool) *filter {
	f.operations = append(f.operations, operation{
		operator:        operator,
		field:           field,
		value:           value,
		isFieldRequired: isFieldRequired,
	})

	return f
}

func (f *filter) ToJSON() (string, error) {
	query, err := f.Build()
	if err != nil {
		return "", err
	}

	b, err := bson.MarshalExtJSON(query, false, false)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
