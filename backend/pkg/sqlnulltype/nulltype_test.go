package sqlnulltype

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullType_String(t *testing.T) {
	str1 := String("")
	assert.Equal(t, false, str1.Valid)
	assert.Equal(t, "", str1.String)

	str2 := String("sinarmas")
	assert.Equal(t, true, str2.Valid)
	assert.Equal(t, "sinarmas", str2.String)
}

func TestNullType_Time(t *testing.T) {
	time1 := time.Time{}
	assert.Equal(t, false, !time1.IsZero())
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", time1.String())

	time2 := time.Date(2023, 11, 28, 21, 22, 56, 647000000, time.UTC)
	assert.Equal(t, true, !time2.IsZero())
	assert.Equal(t, "2023-11-28 21:22:56.647 +0000 UTC", time2.String())
}

func TestNullType_Int64(t *testing.T) {
	val1 := Int64(0)
	assert.Equal(t, false, val1.Valid)
	assert.Equal(t, int64(0), val1.Int64)

	val2 := Int64(2)
	assert.Equal(t, true, val2.Valid)
	assert.Equal(t, int64(2), val2.Int64)
}
