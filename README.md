# result [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/henrylee2cn/result)
Go-generics result module inspired by rust.

Avoid ifelse, handle result with chain methods, will you choose her?

## Go Version

go≥1.18

## Example

```go
func ExampleAndThen() {
	var divide = func(i, j float32) Result[float32] {
		if j == 0 {
			return Err[float32]("j can not be 0")
		}
		return Ok(i / j)
	}
	var ret float32 = divide(1, 2).AndThen(func(i float32) Result[float32] {
		return Ok(i * 10)
	}).Unwrap()
	fmt.Println(ret)
	// Output:
	// 5
}
```
