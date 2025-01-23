package kyte

import "go.mongodb.org/mongo-driver/bson"

const (
	push         = "$push"
	count        = "$count"
	first        = "$first"
	last         = "$last"
	max          = "$max"
	min          = "$min"
	percentile   = "$percentile"
	stdDevPop    = "$stdDevPop"
	stdDevSamp   = "$stdDevSamp"
	sum          = "$sum"
	mergeObjects = "$mergeObjects"

	// TODO
	// avg          = "$avg"
	// addToSet     = "$addToSet"
	// bottom       = "$bottom"
	// bottomN      = "$bottomN"
	// lastN        = "$lastN"
	// maxN         = "$maxN"
	// median       = "$median"
	// minN         = "$minN"
	// top          = "$top"
	// topN         = "$topN"
)

type accumulator struct {
	accumulators []bson.E
}

func Accumulate() *accumulator {
	return &accumulator{}
}

// TODO can expression be more strict? rather than bson.M
// Like Push("field",Expression.Field("name").Field("surname").Field("age"))
func (a *accumulator) Push(field string, expression bson.M) *accumulator {
	a.accumulators = append(a.accumulators, bson.E{Key: field, Value: bson.M{push: expression}})
	return a
}

func (a *accumulator) Count(field string) *accumulator {
	a.accumulators = append(a.accumulators, bson.E{Key: field, Value: bson.M{count: bson.M{}}})
	return a
}

func (a *accumulator) Build() []bson.E {
	return a.accumulators
}
