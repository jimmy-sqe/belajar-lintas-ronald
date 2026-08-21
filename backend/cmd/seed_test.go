package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCollectionFromFilename(t *testing.T) {
	cases := map[string]string{
		"db/seeds/mongodb/001_items.json":        "items",
		"db/seeds/mongodb/042_orders.json":       "orders",
		"db/seeds/mongodb/items.json":            "items",
		"db/seeds/mongodb/abc_users.json":        "abc_users",
		"db/seeds/mongodb/001_complex_name.json": "complex_name",
	}
	for in, want := range cases {
		assert.Equalf(t, want, collectionFromFilename(in), "input=%q", in)
	}
}

func TestConvertDateStrings_ParsesKnownFields(t *testing.T) {
	doc := bson.M{
		"_id":        "abc",
		"name":       "Foo",
		"created_at": "2026-01-01T00:00:00Z",
		"updated_at": "2026-01-02T12:34:56Z",
		"unrelated":  "2026-01-01T00:00:00Z",
	}
	convertDateStrings(doc)

	_, ok := doc["created_at"].(time.Time)
	assert.True(t, ok, "created_at must be time.Time")
	_, ok = doc["updated_at"].(time.Time)
	assert.True(t, ok, "updated_at must be time.Time")
	assert.Equal(t, "2026-01-01T00:00:00Z", doc["unrelated"], "unrelated string field unchanged")
}

func TestConvertDateStrings_LeavesUnparseableAlone(t *testing.T) {
	doc := bson.M{"created_at": "not-a-date"}
	convertDateStrings(doc)
	assert.Equal(t, "not-a-date", doc["created_at"])
}
