package collut

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	t.Run("numeric identity", func(t *testing.T) {
		assert.Equal(t, 1, Identity(1))
		assert.Equal(t, 100, Identity(100))
		assert.Equal(t, 0, Identity(0))
	})

	t.Run("string identity", func(t *testing.T) {
		assert.Equal(t, "hello", Identity("hello"))
	})
}

func TestFuncCompose(t *testing.T) {
	t.Run("apply identity", func(t *testing.T) {
		assert.Equal(t, "123", FuncCompose(Identity[string], Identity[string])("123"))
		assert.Equal(t, "123", FuncCompose(strconv.Itoa, Identity[string])(123))
		assert.Equal(t, "123", FuncCompose(Identity[int], strconv.Itoa)(123))
	})

	t.Run("apply substring", func(t *testing.T) {
		s := assert.New(t)
		substr4 := func(s string) string {
			return s[:4]
		}
		appendWorld := func(s string) string {
			return s + " world"
		}
		s.Equal("HELL", FuncCompose(substr4, strings.ToUpper)("hello"))
		s.Equal("HELL", FuncCompose(strings.ToUpper, substr4)("hello"))
		s.Equal("hell", FuncCompose(appendWorld, substr4)("hello"))
		s.Equal("hell world", FuncCompose(substr4, appendWorld)("hello"))
	})
}

func TestZero(t *testing.T) {
	t.Run("for int it should return 0", func(t *testing.T) {
		assert.Equal(t, 0, Zero[int]())
	})

	t.Run("for string it should return empty string", func(t *testing.T) {
		assert.Empty(t, Zero[string]())
	})

	t.Run("for array it should return empty array", func(t *testing.T) {
		assert.Empty(t, Zero[[]any]())
	})
}

func TestOpsApplyAll(t *testing.T) {
	type testStruct struct {
		value int
	}
	type opFn = func(ts *testStruct)

	incOnce := func(ts *testStruct) {
		ts.value++
	}

	tests := []struct {
		name   string
		ops    []opFn
		expVal int
	}{
		{
			name:   "empty operations",
			ops:    []opFn{},
			expVal: 0,
		},
		{
			name: "single inc operation",
			ops: []opFn{
				incOnce,
			},
			expVal: 1,
		},
		{
			name: "two inc operations",
			ops: []opFn{
				incOnce,
				incOnce,
			},
			expVal: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := OpsApplyAll[testStruct](testStruct{}, tt.ops...)
			assert.Equal(t, tt.expVal, res.value)
		})
	}
}
