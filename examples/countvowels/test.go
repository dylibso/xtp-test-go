package main

import (
	"encoding/json"
	"fmt"
	"math"

	xtptest "github.com/dylibso/xtp-test-go"
)

//go:export test
func test() int32 {
	xtptest.AssertGt("gt test", 100, 1)
	xtptest.AssertLt("lt test", 1, 100)
	xtptest.AssertGte("gte test", math.Pi, 3.14)
	xtptest.AssertLte("lte test", 3.14, math.Pi)

	output := xtptest.CallString("count_vowels", []byte("hello"))
	xtptest.AssertNe("we got some output", output, "")

	xtptest.Group("check how fast the function performs", func() {
		sec := xtptest.TimeSeconds("count_vowels", []byte("hello"))
		xtptest.AssertLt("it should be fast", sec, 0.1)
		ns := xtptest.TimeNanos("count_vowels", []byte("hello"))
		xtptest.AssertLt("it should be really fast", ns, 300000)
	})

	xtptest.Group("check that count_vowels maintains state", func() {
		accumTotal := 0
		for i := 0; i < 10; i++ {
			c := fromJson(xtptest.CallBytes("count_vowels", []byte("this is a test")))
			accumTotal += c.Count
			xtptest.AssertEq(fmt.Sprintf("total should be incremented to: %d", accumTotal), c.Total, accumTotal)
		}
	})

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
