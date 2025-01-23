package kyte

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	match   = "$match"
	group   = "$group"
	project = "$project"
	sort    = "$sort"
	limit   = "$limit"
	skip    = "$skip"
	unwind  = "$unwind"
	lookup  = "$lookup"
	facet   = "$facet"

	// TODO
	// addField = "$addFields"
	// bucket = "$bucket"
	// bucketAuto = "$bucketAuto"
	// chageStream = "$changeStream"
	// changeStreamSplitLargeEvent = "$changeStreamSplitLargeEvent"
	// collStats = "$collStats"
	// count = "$count"
	// densify = "$densify"
	// documents = "$documents"
	// fill = "$fill"
	// geoNear = "$geoNear"
	// graphLookup = "$graphLookup"
	// indexStats = "$indexStats"
	// listSampledQueries = "$listSampledQueries"
	// listSearchIndexes = "$listSearchIndexes"
	// listSessions = "$listSessions"
	// merge = "$merge"
	// out = "$out"
	// planCacheStats = "$planCacheStats"
	// redact = "$redact"
	// replaceRoot = "$replaceRoot"
	// replaceWith = "$replaceWith"
	// sample = "$sample"
	// search = "$search"
	// searchMeta = "$searchMeta"
	// set = "$set"
	// setWindowFields = "$setWindowFields"
	// sortByCount = "$sortByCount"
	// unionWith = "$unionWith"
	// unset = "$unset"
)

type aggregate struct {
	kyte       *kyte
	pipeline   mongo.Pipeline
	operations []operation
}

/*
Aggregate creates a new aggregate instance.
*/
func Aggregate(opts ...OptionFunc) *aggregate {
	options := &Options{validateField: true}
	for _, opt := range opts {
		opt(options)
	}

	kyte := newKyte(options.source, options.validateField)

	return &aggregate{
		kyte:     kyte,
		pipeline: mongo.Pipeline{},
	}
}

/*
Match adds a [$match] operation to the aggregate pipeline.

	Aggregate().
		Match(
				Filter().
					Exists("age", true).
					GreaterThan("age", 10).
					LessThan("age", 20).
					Type("name", bson.TypeString),
		).
		Build()

[$match]: https://www.mongodb.com/docs/manual/reference/operator/aggregation/match/#mongodb-pipeline-pipe.-match
*/
func (a *aggregate) Match(filter *filter) *aggregate {
	if a.kyte.source != nil {
		filter.kyte.checkField = a.kyte.checkField
		filter.kyte.setSourceAndPrepareFields(a.kyte.source)
	}

	query, err := filter.Build()
	if err != nil {
		a.kyte.setError(err)
		return a
	}

	a.pipeline = append(a.pipeline, bson.D{{Key: match, Value: query}})
	return a
}

func (a *aggregate) Group(_id any, acc *accumulator) *aggregate {
	accBuild := acc.Build()
	if len(accBuild) == 0 {
		return a
	}

	values := bson.M{
		"_id": _id,
	}

	for _, v := range accBuild {
		values[v.Key] = v.Value
	}

	a.pipeline = append(a.pipeline, bson.D{{Key: group, Value: values}})
	return a
}

func (a *aggregate) Build() (mongo.Pipeline, error) {
	if a.kyte.err != nil {
		return nil, a.kyte.err
	}

	return a.pipeline, nil
}
