package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLOffset(t *testing.T) {
	tests := []struct {
		name                          string
		page, perPage, expectedResult uint64
	}{
		{
			name:           "offset first page",
			page:           1,
			perPage:        10,
			expectedResult: 0,
		},
		{
			name:           "offset nth page",
			page:           5,
			perPage:        8,
			expectedResult: 32,
		},
		{
			name:           "offset zero page",
			page:           0,
			perPage:        8,
			expectedResult: 0,
		},
	}

	for _, tt := range tests {
		result := SQLOffset(tt.page, tt.perPage)
		assert.Equal(t, tt.expectedResult, result)
	}
}

func TestTotalPage(t *testing.T) {
	tests := []struct {
		name                              string
		totalRow, perPage, expectedResult uint64
	}{
		{
			name:           "0 row",
			totalRow:       0,
			perPage:        10,
			expectedResult: 0,
		},
		{
			name:           "1 row",
			totalRow:       1,
			perPage:        10,
			expectedResult: 1,
		},
		{
			name:           "1 page",
			totalRow:       10,
			perPage:        10,
			expectedResult: 1,
		},
		{
			name:           "multiple page",
			totalRow:       34,
			perPage:        8,
			expectedResult: 5,
		},
	}
	for _, tt := range tests {
		result := TotalPage(tt.totalRow, tt.perPage)
		assert.Equal(t, tt.expectedResult, result)
	}
}

func TestBuild(t *testing.T) {
	tests := []struct {
		name      string
		page      uint64
		pageSize  uint64
		totalData uint64
		want      Page
	}{
		{
			name:      "first page, exact fit",
			page:      1,
			pageSize:  10,
			totalData: 10,
			want:      Page{Page: 1, PageSize: 10, TotalData: 10, TotalPage: 1},
		},
		{
			name:      "multi page, not divisible",
			page:      2,
			pageSize:  8,
			totalData: 34,
			want:      Page{Page: 2, PageSize: 8, TotalData: 34, TotalPage: 5},
		},
		{
			name:      "zero total",
			page:      1,
			pageSize:  10,
			totalData: 0,
			want:      Page{Page: 1, PageSize: 10, TotalData: 0, TotalPage: 0},
		},
		{
			name:      "page zero is normalised to 1",
			page:      0,
			pageSize:  10,
			totalData: 5,
			want:      Page{Page: 1, PageSize: 10, TotalData: 5, TotalPage: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Build(tt.page, tt.pageSize, tt.totalData))
		})
	}
}
