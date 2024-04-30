package main

import (
	"encoding/json"
	"fmt"

	xtptest "github.com/dylibso/xtp-test-go"
)

//go:export test
func test() int32 {
	output := xtptest.CallString("count_vowels", []byte("hello"))
	xtptest.AssertNe("we got some output", output, "")

	xtptest.Group("check how fast the function performs", func() {
		sec := xtptest.TimeSeconds("count_vowels", []byte("hello"))
		xtptest.Assert("it should be fast", sec < 0.1, fmt.Sprintf("Expected %f < 0.1", sec))
		ns := xtptest.TimeNanos("count_vowels", []byte("hello"))
		xtptest.Assert("it should be really fast", ns < 100000, fmt.Sprintf("Expected %d < 100000", ns))
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
