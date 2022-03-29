package result

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleAndThan() {
	var divide = func(i, j float32) Result[float32] {
		if j == 0 {
			return Err[float32]("j can not be 0")
		}
		return Ok(i / j)
	}
	var ret float32 = divide(1, 2).AndThan(func(i float32) Result[float32] {
		return Ok(i * 10)
	}).Unwrap()
	fmt.Println(ret)
	// Output:
	// 5
}

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
