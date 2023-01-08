package collut

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceMap(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		assert.Empty(t, SliceMap([]int{}, func(i int) int { return 0 }))
	})

	t.Run("slices with int to int", func(t *testing.T) {
		s := assert.New(t)
		s.Equal([]int{1, 2, 3}, SliceMap([]int{1, 2, 3}, Identity[int]))
		s.Equal([]int{0, 0}, SliceMap([]int{1, 2}, func(_ int) int { return 0 }))
		s.Equal([]int{2, 4}, SliceMap([]int{1, 2}, func(i int) int { return i * 2 }))
	})

	t.Run("slices with int to string", func(t *testing.T) {
		s := assert.New(t)
		s.Equal([]string{"1", "2", "3"}, SliceMap([]int{1, 2, 3}, strconv.Itoa))
	})

	t.Run("slice of struct, extract id", func(t *testing.T) {
		type ti struct {
			id   string
			name string
		}
		items := []ti{
			{id: "id1", name: "name1"},
			{id: "id2", name: "name2"},
		}
		res := SliceMap(items, func(i ti) string { return i.id })
		assert.Equal(t, []string{"id1", "id2"}, res)
	})
}

func TestSliceFilter(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		s := assert.New(t)
		s.Empty(SliceFilter([]int{}, func(i int) bool { return true }))
	})

	t.Run("slice of integers", func(t *testing.T) {
		s := assert.New(t)
		items := []int{1, 2, 3, 4}
		s.Equal(items, SliceFilter(items, func(i int) bool { return true }))
		s.Empty(SliceFilter(items, func(i int) bool { return false }))
		s.Equal([]int{2, 4}, SliceFilter(items, func(i int) bool { return i%2 == 0 }))
	})
}

func TestSliceContains(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		s := assert.New(t)
		s.False(SliceContains([]int{}, func(i int) bool { return true }))
	})

	t.Run("non-empty int slice", func(t *testing.T) {
		tests := []struct {
			val int
			exp bool
		}{
			{0, false},
			{4, false},
			{1, true},
			{2, true},
			{3, true},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("for number: %d", test.val), func(t *testing.T) {
				s := assert.New(t)
				numbers := []int{1, 2, 3}
				if test.exp {
					s.True(SliceContains(numbers, func(i int) bool { return i == test.val }))
				} else {
					s.False(SliceContains(numbers, func(i int) bool { return i == test.val }))
				}
			})
		}
	})
}

func TestSliceFoldl(t *testing.T) {
	intSum := func(acc, i int) int { return acc + i }
	items := []int{1, 2, 3}
	t.Run("empty slice", func(t *testing.T) {
		s := assert.New(t)
		s.Equal(0, SliceFoldl(0, []int{}, intSum))
	})

	t.Run("sum of integers", func(t *testing.T) {
		s := assert.New(t)
		s.Equal(6, SliceFoldl(0, items, intSum))
	})
	t.Run("appender", func(t *testing.T) {
		s := assert.New(t)
		res := SliceFoldl([]int{}, items, func(acc []int, i int) []int {
			return append(acc, i)
		})
		s.Equal([]int{1, 2, 3}, res)
	})
}

func TestSliceFoldr(t *testing.T) {
	intSum := func(acc, i int) int { return acc + i }
	items := []int{1, 2, 3}

	t.Run("empty slice", func(t *testing.T) {
		s := assert.New(t)
		s.Equal(0, SliceFoldr(0, []int{}, intSum))
	})

	t.Run("sum of integers", func(t *testing.T) {
		s := assert.New(t)
		s.Equal(6, SliceFoldr(0, items, intSum))
	})

	t.Run("appender", func(t *testing.T) {
		s := assert.New(t)
		res := SliceFoldr([]int{}, items, func(acc []int, i int) []int {
			return append(acc, i)
		})
		s.Equal([]int{3, 2, 1}, res)
	})
}

func TestSliceLengthLimit(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		limit int
		exp   []int
	}{
		{
			name:  "empty slice, zero limit",
			input: []int{},
			limit: 0,
			exp:   []int{},
		},
		{
			name:  "empty slice, non-zero limit",
			input: []int{},
			limit: 10,
			exp:   []int{},
		},
		{
			name:  "non-empty slice, non-zero limit (n-1)",
			input: []int{1, 2, 3},
			limit: 2,
			exp:   []int{1, 2},
		},
		{
			name:  "non-empty slice, zero limit",
			input: []int{1, 2, 3},
			limit: 0,
			exp:   []int{},
		},
		{
			name:  "non-empty slice, large limit",
			input: []int{1, 2, 3},
			limit: 10,
			exp:   []int{1, 2, 3},
		},
		{
			name:  "non-empty slice, exact limit",
			input: []int{1, 2, 3},
			limit: 3,
			exp:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := assert.New(t)
			s.Equal(tt.exp, SliceLengthLimit(tt.input, tt.limit))
		})
	}
}
