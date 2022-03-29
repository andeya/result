package result

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var w = Wrap[int](strconv.Atoi("s"))
	assert.False(t, w.IsOk())
	assert.True(t, w.IsErr())
	assert.Nil(t, w.Ok())

	var w2 = Wrap[any](strconv.Atoi("-1"))
	assert.True(t, w2.IsOk())
	assert.False(t, w2.IsErr())
	assert.Equal(t, -1, *w2.Ok())
}

func TestResult_Map(t *testing.T) {
	var isMyNum = func(s string, search int) Result[bool] {
		return Map(Wrap(strconv.Atoi(s)), func(x int) bool { return x == search })
	}
	assert.Equal(t, Ok[bool](true), isMyNum("1", 1))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}
