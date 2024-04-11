# xtp-test

A Go test framework for [xtp](https://getxtp.com) / [Extism](https://extism.org)
plugins.

## Example

```go
package main

import (
	"encoding/json"
	"fmt"

	xtptest "github.com/dylibso/xtp-test-go"
)

// You _must_ export a single `test` function in order for the test runner to call something
//
//go:export test
func test() int32 {
	// call the tested plugin's "count_vowels" function, passing it some data
	output := xtptest.CallString("count_vowels", []byte("hello"))
	// check that the output is as expected
	xtptest.AssertNe("we got some output", output, "")

	// create a named group of tests. NOTE: plugin state is reset before and after the group runs.
	xtptest.Group("check how fast the function performs", func() {
		// check the amount of time in seconds or nanoseconds spent in the plugin function.
		sec := xtptest.TimeSeconds("count_vowels", []byte("hello"))
		xtptest.Assert("it should be fast", sec < 0.1)
		ns := xtptest.TimeNanos("count_vowels", []byte("hello"))
		xtptest.Assert("it should be really fast", ns < 100000)
	})

	xtptest.Group("check that count_vowels maintains state", func() {
		accumTotal := 0
		for i := 0; i < 10; i++ {
			c := fromJson(xtptest.CallBytes("count_vowels", []byte("this is a test")))
			accumTotal += c.Count
			xtptest.AssertEq(fmt.Sprintf("total should be incremented to: %d", accumTotal), c.Total, accumTotal)
		}
	})

	// as this is an Extism plugin, return a status/error code
	return 0
}

func fromJson(data []byte) CountVowels {
	var res CountVowels
	_ = json.Unmarshal(data, &res)
	return res
}

type CountVowels struct {
	Vowels string `json:"vowels"`
	Total  int    `json:"total"`
	Count  int    `json:"count"`
}

func main() {}
```

## API Docs

Please see the [**`godoc`**](https://pkg.go.dev/github.com/dylibso/xtp-test-go)
documentation page for full details.

## Usage

**1. Create a Go project using the XTP Test crate**

```sh
mkdir plugintest && cd plugintest
go mod init <your module>
go get github.com/dylibso/xtp-test-go
```

**2. Write your test in Go**

```go
// test.go
package main

import (
	"encoding/json"
	"fmt"

	xtptest "github.com/dylibso/xtp-test-go"
)

//go:export test
func test() int32 {
    // call the tested plugin's "count_vowels" function, passing it some data
	output := xtptest.CallString("count_vowels", []byte("hello"))
    // check that the output is as expected
	xtptest.AssertNe("we got some output", output, "")
    // ...
```

**3. Compile your test to .wasm:**

Ensure you have
[`tinygo` installed](https://tinygo.org/getting-started/install/). Eventually
you can use the `go` compiler, but the test runner _must_ find an exported
`test` function from the test plugin, and `go` cannot currently export any
function other than `_start`.

```sh
tinygo build -o test.wasm -target wasi test.go
```

**4. Run the test against your plugin:** Once you have your test code as a
`.wasm` module, you can run the test against your plugin using the `xtp` CLI:

### Install `xtp`

```sh
curl https://static.dylibso.com/cli/install.sh | sudo sh
```

### Run the test suite

```sh
xtp plugin test ./plugin-*.wasm --with test.wasm --host host.wasm
#               ^^^^^^^^^^^^^^^        ^^^^^^^^^        ^^^^^^^^^
#               your plugin(s)         test to run      optional mock host functions
```

**Note:** The optional mock host functions must be implemented as Extism
plugins, whose exported functions match the host function signature imported by
the plugins being tested.

## Need Help?

Please reach out via the
[`#xtp` channel on Discord](https://discord.com/channels/1011124058408112148/1220464672784908358).
