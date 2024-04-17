package main

import (
	"fmt"
	"strings"

	xtptest "github.com/dylibso/xtp-test-go"
)

//go:export test
func test() int32 {
	output := xtptest.CallString("run", nil)
	xtptest.AssertEq("initial call to 'run' returns the expected value", output, "value")

	xtptest.Group("multiple kv read/write calls produce correct state", func() {
		for i := 0; i < 10; i++ {
			output := xtptest.CallString("run", nil)
			expected := strings.Repeat("value", i+1)
			msg := fmt.Sprintf("repeat call to 'run' returns the correct value: %s", expected)
			xtptest.AssertEq(msg, output, expected)
		}
	})

	return 0
}

func main() {}
