// Package integration_test executes Kyte Queries against a real MongoDB
// server started with testcontainers-go and asserts which documents match.
//
// The tests are skipped unless KYTE_INTEGRATION_TEST=true. When it is set,
// Docker must be reachable or the run fails. KYTE_MONGO_IMAGE selects the
// MongoDB image (default mongo:8.0).
package integration_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	enableEnv    = "KYTE_INTEGRATION_TEST"
	imageEnv     = "KYTE_MONGO_IMAGE"
	defaultImage = "mongo:8.0"
	mongoPort    = "27017/tcp"
)

// db is nil when integration tests are disabled.
var db *mongo.Database

// buildable is satisfied by the Filter returned from kyte.Filter(), whose type is unexported.
type buildable interface {
	Build() (bson.D, error)
}

func TestMain(m *testing.M) {
	if os.Getenv(enableEnv) != "true" {
		os.Exit(m.Run())
	}

	os.Exit(runWithMongo(m))
}

func runWithMongo(m *testing.M) int {
	ctx := context.Background()

	image := os.Getenv(imageEnv)
	if image == "" {
		image = defaultImage
	}

	container, err := testcontainers.Run(ctx, image,
		testcontainers.WithExposedPorts(mongoPort),
		testcontainers.WithWaitStrategy(wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort(mongoPort),
		)),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			fmt.Fprintf(os.Stderr, "terminate %s container: %v\n", image, err)
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "start %s container (is Docker running?): %v\n", image, err)
		return 1
	}

	uri, err := container.PortEndpoint(ctx, mongoPort, "mongodb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve %s endpoint: %v\n", image, err)
		return 1
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to %s: %v\n", uri, err)
		return 1
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Fprintf(os.Stderr, "ping %s: %v\n", uri, err)
		return 1
	}

	db = client.Database("kyte")
	return m.Run()
}

var invalidCollectionChars = regexp.MustCompile(`[^A-Za-z0-9_]`)

// seed inserts the fixture users into a collection owned by t and returns it.
// It skips t when integration tests are disabled.
func seed(t *testing.T) *mongo.Collection {
	t.Helper()

	if db == nil {
		t.Skipf("set %s=true to run integration tests against a MongoDB container", enableEnv)
	}

	coll := db.Collection(invalidCollectionChars.ReplaceAllString(t.Name(), "_"))

	docs := make([]any, len(users))
	for i, u := range users {
		docs[i] = u
	}

	if _, err := coll.InsertMany(t.Context(), docs); err != nil {
		t.Fatalf("seed %s: %v", coll.Name(), err)
	}

	return coll
}

// assertMatches runs the Query built by f with Find and checks that exactly the
// users with the wanted _ids match, in any order.
//
// want must be a non-empty proper subset of the fixtures: a broken Query that
// matches nothing or everything would otherwise pass unnoticed.
func assertMatches(t *testing.T, coll *mongo.Collection, f buildable, want ...int) {
	t.Helper()

	if len(want) == 0 || len(want) >= len(users) {
		t.Fatalf("want must be a non-empty proper subset of the %d fixture users, got %v", len(users), want)
	}

	query, err := f.Build()
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	opts := options.Find().
		SetProjection(bson.D{{Key: "_id", Value: 1}}).
		SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := coll.Find(t.Context(), query, opts)
	if err != nil {
		t.Fatalf("Find(%s) returned error: %v", toJSON(query), err)
	}

	var matched []struct {
		ID int `bson:"_id"`
	}
	if err := cursor.All(t.Context(), &matched); err != nil {
		t.Fatalf("decode Find(%s) results: %v", toJSON(query), err)
	}

	got := make([]int, len(matched))
	for i, m := range matched {
		got[i] = m.ID
	}

	want = slices.Clone(want)
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Errorf("Find(%s) matched _ids %v, want %v", toJSON(query), got, want)
	}
}

func toJSON(query bson.D) string {
	b, err := bson.MarshalExtJSON(query, false, false)
	if err != nil {
		return fmt.Sprintf("%v", query)
	}
	return string(b)
}
