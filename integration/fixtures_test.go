package integration_test

// User is the Source for every integration test. Its shape covers each way a
// field pointer resolves to a bson key: top-level, nested struct, struct two
// levels deep, pointer to struct and slice of structs.
type User struct {
	ID       int      `bson:"_id"`
	Name     string   `bson:"name"`
	Age      int      `bson:"age"`
	Tenant   string   `bson:"tenant"`
	Tags     []string `bson:"tags"`
	Address  Address  `bson:"address"`
	Profile  *Profile `bson:"profile,omitempty"`
	Orders   []Order  `bson:"orders"`
	Nickname *string  `bson:"nickname,omitempty"`
	// Score holds a different BSON type per user (int32, double, string, null).
	Score any `bson:"score"`
}

type Address struct {
	City    string `bson:"city"`
	Country string `bson:"country"`
	Geo     Geo    `bson:"geo"`
}

type Geo struct {
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
}

type Profile struct {
	Bio string `bson:"bio"`
}

type Order struct {
	Item     string `bson:"item"`
	Quantity int    `bson:"quantity"`
}

// newSource returns a Source whose pointer and slice fields are populated, so
// that pointers into them (&u.Profile.Bio, &u.Orders[0].Item) can be taken.
func newSource() *User {
	return &User{Profile: &Profile{}, Orders: []Order{{}}}
}

func ptr[T any](v T) *T {
	return &v
}

// users is the fixture set seeded into every test collection.
//
//	_id name   age tenant tags               city/country  profile orders            nickname score
//	1   Alice  25  t1     [go mongo]         Istanbul/TR   gopher  [book×2]          ali      int32 90
//	2   alex   30  t1     [go]               Ankara/TR     -       [pen×10 book×1]   -        double 85.5
//	3   Bob    35  t2     [mongo db go]      Berlin/DE     dba     []                bobby    string "A"
//	4   Carol  40  t2     []                 Munich/DE     -       [laptop×1]        -        null
//	5   Dave   45  t1     [db]               Izmir/TR      ops     [pen×3]           d        int32 70
//	6   Eve    50  t2     [go db]            Paris/FR      -       []                -        string "B"
var users = []User{
	{
		ID: 1, Name: "Alice", Age: 25, Tenant: "t1",
		Tags:     []string{"go", "mongo"},
		Address:  Address{City: "Istanbul", Country: "TR", Geo: Geo{Lat: 41.0, Lng: 28.9}},
		Profile:  &Profile{Bio: "gopher"},
		Orders:   []Order{{Item: "book", Quantity: 2}},
		Nickname: ptr("ali"),
		Score:    int32(90),
	},
	{
		ID: 2, Name: "alex", Age: 30, Tenant: "t1",
		Tags:    []string{"go"},
		Address: Address{City: "Ankara", Country: "TR", Geo: Geo{Lat: 39.9, Lng: 32.8}},
		Orders:  []Order{{Item: "pen", Quantity: 10}, {Item: "book", Quantity: 1}},
		Score:   85.5,
	},
	{
		ID: 3, Name: "Bob", Age: 35, Tenant: "t2",
		Tags:     []string{"mongo", "db", "go"},
		Address:  Address{City: "Berlin", Country: "DE", Geo: Geo{Lat: 52.5, Lng: 13.4}},
		Profile:  &Profile{Bio: "dba"},
		Orders:   []Order{},
		Nickname: ptr("bobby"),
		Score:    "A",
	},
	{
		ID: 4, Name: "Carol", Age: 40, Tenant: "t2",
		Tags:    []string{},
		Address: Address{City: "Munich", Country: "DE", Geo: Geo{Lat: 48.1, Lng: 11.6}},
		Orders:  []Order{{Item: "laptop", Quantity: 1}},
		Score:   nil,
	},
	{
		ID: 5, Name: "Dave", Age: 45, Tenant: "t1",
		Tags:     []string{"db"},
		Address:  Address{City: "Izmir", Country: "TR", Geo: Geo{Lat: 38.4, Lng: 27.1}},
		Profile:  &Profile{Bio: "ops"},
		Orders:   []Order{{Item: "pen", Quantity: 3}},
		Nickname: ptr("d"),
		Score:    int32(70),
	},
	{
		ID: 6, Name: "Eve", Age: 50, Tenant: "t2",
		Tags:    []string{"go", "db"},
		Address: Address{City: "Paris", Country: "FR", Geo: Geo{Lat: 48.9, Lng: 2.4}},
		Orders:  []Order{},
		Score:   "B",
	},
}
