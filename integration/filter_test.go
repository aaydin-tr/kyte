package integration_test

import (
	"regexp"
	"testing"

	"github.com/aaydin-tr/kyte"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type matchCase struct {
	name   string
	filter buildable
	want   []int
}

// runMatchCases seeds a collection for t and asserts each case against it in
// parallel subtests. The cases only read, so they share the seeded collection.
func runMatchCases(t *testing.T, cases []matchCase) {
	t.Helper()

	coll := seed(t)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			assertMatches(t, coll, c.filter, c.want...)
		})
	}
}

func TestFilter_Equal(t *testing.T) {
	t.Parallel()

	u := newSource()
	name := "Bob"

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().Equal("name", "Bob"), []int{3}},
		{"source field", kyte.Filter(kyte.Source(u)).Equal(&u.Name, "Bob"), []int{3}},
		{"pointer value", kyte.Filter().Equal("name", &name), []int{3}},
		{"array element", kyte.Filter().Equal("tags", "db"), []int{3, 5, 6}},
		{"dotted string field", kyte.Filter().Equal("address.geo.lat", 52.5), []int{3}},
		{"nested struct field", kyte.Filter(kyte.Source(u)).Equal(&u.Address.City, "Berlin"), []int{3}},
		{"pointer to struct field", kyte.Filter(kyte.Source(u)).Equal(&u.Profile.Bio, "dba"), []int{3}},
		{"slice of structs field", kyte.Filter(kyte.Source(u)).Equal(&u.Orders[0].Item, "pen"), []int{2, 5}},
	})

	t.Run("struct nested two levels deep", func(t *testing.T) {
		t.Skip(`known bug: &u.Address.Geo.Lat resolves to "geo..lat" instead of "address.geo.lat"`)
		assertMatches(t, seed(t), kyte.Filter(kyte.Source(u)).Equal(&u.Address.Geo.Lat, 52.5), 3)
	})
}

func TestFilter_NotEqual(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().NotEqual("tenant", "t1"), []int{3, 4, 6}},
		{"source field", kyte.Filter(kyte.Source(u)).NotEqual(&u.Tenant, "t2"), []int{1, 2, 5}},
		{"matches missing field", kyte.Filter().NotEqual("nickname", "ali"), []int{2, 3, 4, 5, 6}},
	})
}

func TestFilter_GreaterThan(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().GreaterThan("age", 40), []int{5, 6}},
		{"source field", kyte.Filter(kyte.Source(u)).GreaterThan(&u.Age, 40), []int{5, 6}},
		{"mixed numeric types", kyte.Filter().GreaterThan("score", 80), []int{1, 2}},
		{"slice of structs field", kyte.Filter(kyte.Source(u)).GreaterThan(&u.Orders[0].Quantity, 5), []int{2}},
	})
}

func TestFilter_GreaterThanOrEqual(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().GreaterThanOrEqual("age", 40), []int{4, 5, 6}},
		{"source field", kyte.Filter(kyte.Source(u)).GreaterThanOrEqual(&u.Age, 40), []int{4, 5, 6}},
	})
}

func TestFilter_LessThan(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().LessThan("age", 30), []int{1}},
		{"source field", kyte.Filter(kyte.Source(u)).LessThan(&u.Age, 30), []int{1}},
	})
}

func TestFilter_LessThanOrEqual(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().LessThanOrEqual("age", 30), []int{1, 2}},
		{"source field", kyte.Filter(kyte.Source(u)).LessThanOrEqual(&u.Age, 30), []int{1, 2}},
	})
}

func TestFilter_In(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"slice value", kyte.Filter().In("address.country", []string{"DE", "FR"}), []int{3, 4, 6}},
		{"scalar value", kyte.Filter().In("name", "Eve"), []int{6}},
		{"source field", kyte.Filter(kyte.Source(u)).In(&u.Address.Country, []string{"FR"}), []int{6}},
		{"array field", kyte.Filter().In("tags", []string{"mongo"}), []int{1, 3}},
	})
}

func TestFilter_NotIn(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"slice value", kyte.Filter().NotIn("address.country", []string{"TR", "DE"}), []int{6}},
		{"scalar value", kyte.Filter().NotIn("tenant", "t1"), []int{3, 4, 6}},
		{"source field", kyte.Filter(kyte.Source(u)).NotIn(&u.Tenant, []string{"t2"}), []int{1, 2, 5}},
	})
}

func TestFilter_Regex(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"without options", kyte.Filter().Regex("name", regexp.MustCompile("^a")), []int{2}},
		{"with options", kyte.Filter().Regex("name", regexp.MustCompile("^a"), "i"), []int{1, 2}},
		{"source field", kyte.Filter(kyte.Source(u)).Regex(&u.Address.City, regexp.MustCompile("^I")), []int{1, 5}},
	})
}

func TestFilter_Exists(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"true", kyte.Filter().Exists("nickname", true), []int{1, 3, 5}},
		{"false", kyte.Filter().Exists("nickname", false), []int{2, 4, 6}},
		{"source field", kyte.Filter(kyte.Source(u)).Exists(&u.Profile, false), []int{2, 4, 6}},
	})
}

func TestFilter_Type(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"single type", kyte.Filter().Type("score", bsontype.String), []int{3, 6}},
		{"multiple types", kyte.Filter().Type("score", bsontype.Int32, bsontype.Double), []int{1, 2, 5}},
		{"null", kyte.Filter().Type("score", bsontype.Null), []int{4}},
		{"source field", kyte.Filter(kyte.Source(u)).Type(&u.Score, bsontype.Int32), []int{1, 5}},
	})
}

func TestFilter_Mod(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().Mod("age", 10, 0), []int{2, 4, 6}},
		{"source field", kyte.Filter(kyte.Source(u)).Mod(&u.Age, 10, 5), []int{1, 3, 5}},
	})
}

func TestFilter_Where(t *testing.T) {
	t.Parallel()

	runMatchCases(t, []matchCase{
		{"array length", kyte.Filter().Where("this.tags.length > 1"), []int{1, 3, 6}},
		{"string function", kyte.Filter().Where("this.name.toLowerCase().startsWith('a')"), []int{1, 2}},
	})
}

func TestFilter_All(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().All("tags", []string{"go", "mongo"}), []int{1, 3}},
		{"source field", kyte.Filter(kyte.Source(u)).All(&u.Tags, []string{"db"}), []int{3, 5, 6}},
	})
}

func TestFilter_Size(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string field", kyte.Filter().Size("tags", 1), []int{2, 5}},
		{"empty array", kyte.Filter().Size("tags", 0), []int{4}},
		{"source field", kyte.Filter(kyte.Source(u)).Size(&u.Orders, 0), []int{3, 6}},
	})
}

func TestFilter_JSONSchema(t *testing.T) {
	t.Parallel()

	runMatchCases(t, []matchCase{
		{"required", kyte.Filter().JSONSchema(bson.M{"required": bson.A{"nickname"}}), []int{1, 3, 5}},
		{"properties", kyte.Filter().JSONSchema(bson.M{
			"properties": bson.M{"age": bson.M{"minimum": 45}},
		}), []int{5, 6}},
		{"required and properties", kyte.Filter().JSONSchema(bson.M{
			"required":   bson.A{"profile"},
			"properties": bson.M{"age": bson.M{"maximum": 30}},
		}), []int{1}},
	})
}

func TestFilter_Raw(t *testing.T) {
	t.Parallel()

	runMatchCases(t, []matchCase{
		{"alone", kyte.Filter().Raw(bson.D{{Key: "address.country", Value: "FR"}}), []int{6}},
		{"with operator", kyte.Filter().Raw(bson.D{{Key: "tenant", Value: "t1"}}).GreaterThan("age", 25), []int{2, 5}},
	})
}

func TestFilter_MultipleOperators(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"string fields", kyte.Filter().Equal("tenant", "t2").LessThan("age", 45), []int{3, 4}},
		{"source fields", kyte.Filter(kyte.Source(u)).
			Equal(&u.Address.Country, "TR").
			Exists(&u.Nickname, true).
			GreaterThan(&u.Age, 30), []int{5}},
	})
}

func TestFilter_And(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"alone", kyte.Filter().And(kyte.Filter().Equal("tenant", "t1").GreaterThan("age", 25)), []int{2, 5}},
		{"with operator", kyte.Filter().Equal("address.country", "TR").And(kyte.Filter().LessThan("age", 30)), []int{1}},
		{"nested nor", kyte.Filter().And(kyte.Filter().
			Equal("tenant", "t2").
			NOR(kyte.Filter().Equal("name", "Eve"))), []int{3, 4}},
		{"source passed to child", kyte.Filter(kyte.Source(u)).And(kyte.Filter().
			Equal(&u.Tenant, "t2").
			GreaterThan(&u.Age, 40)), []int{6}},
	})
}

func TestFilter_Or(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"alone", kyte.Filter().Or(kyte.Filter().Equal("name", "Bob").Equal("name", "Eve")), []int{3, 6}},
		{"with operator", kyte.Filter().Equal("tenant", "t2").Or(kyte.Filter().LessThan("age", 40).GreaterThan("age", 45)), []int{3, 6}},
		{"nested and", kyte.Filter().Or(kyte.Filter().
			Equal("name", "Alice").
			And(kyte.Filter().Equal("tenant", "t2").GreaterThan("age", 45))), []int{1, 6}},
		{"source passed to child", kyte.Filter(kyte.Source(u)).Or(kyte.Filter().
			Equal(&u.Address.City, "Berlin").
			Equal(&u.Address.City, "Paris")), []int{3, 6}},
	})
}

func TestFilter_NOR(t *testing.T) {
	t.Parallel()

	u := newSource()

	runMatchCases(t, []matchCase{
		{"alone", kyte.Filter().NOR(kyte.Filter().Equal("tenant", "t1").Equal("address.country", "DE")), []int{6}},
		{"with operator", kyte.Filter().Equal("tenant", "t1").NOR(kyte.Filter().LessThan("age", 30)), []int{2, 5}},
		{"nested or", kyte.Filter().NOR(kyte.Filter().
			Or(kyte.Filter().Equal("tenant", "t1").Equal("address.country", "FR"))), []int{3, 4}},
		{"source passed to child", kyte.Filter(kyte.Source(u)).NOR(kyte.Filter().
			Equal(&u.Tenant, "t1")), []int{3, 4, 6}},
	})
}
