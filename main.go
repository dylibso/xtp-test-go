package xtptest

import (
	"cmp"
	"fmt"

	"github.com/extism/go-pdk"
)

//go:wasmimport xtp:test/harness call
func call(name uint64, input uint64) uint64

//go:wasmimport xtp:test/harness time
func time(name uint64, input uint64) uint64

//go:wasmimport xtp:test/harness assert
func assert(name uint64, value uint64, reason uint64)

//go:wasmimport xtp:test/harness reset
func reset()

//go:wasmimport xtp:test/harness group
func group(name uint64)

//go:wasmimport xtp:test/harness mock_input
func mock_input() uint64

// Call a function from the Extism plugin being tested, passing input and returning its output Memory.
func Call(funcName string, input []byte) pdk.Memory {
	funcMem := pdk.AllocateString(funcName)
	inputMem := pdk.AllocateBytes(input)

	outputPtr := call(funcMem.Offset(), inputMem.Offset())
	funcMem.Free()
	inputMem.Free()

	return pdk.FindMemory(outputPtr)

}

// Call a function from the Extism plugin being tested, passing input and returning its output as []byte.
func CallBytes(funcName string, input []byte) []byte {
	outputMem := Call(funcName, input)
	return outputMem.ReadBytes()
}

// Call a function from the Extism plugin being tested, passing input and returning its output as a string.
func CallString(funcName string, input []byte) string {
	return string(CallBytes(funcName, input))
}

// Call a function from the Extism plugin being tested, passing input and returning the time in nanoseconds spent in the fuction.
func TimeNanos(funcName string, input []byte) uint64 {
	funcMem := pdk.AllocateString(funcName)
	inputMem := pdk.AllocateBytes(input)

	outputPtr := time(funcMem.Offset(), inputMem.Offset())
	funcMem.Free()
	inputMem.Free()

	return outputPtr
}

// Call a function from the Extism plugin being tested, passing input and returning the time in seconds spent in the fuction.
func TimeSeconds(funcName string, input []byte) float64 {
	return float64(TimeNanos(funcName, input)) / 1e9
}

// Assert that the `outcome` is true, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func Assert(msg string, outcome bool, reason string) {
	msgMem := pdk.AllocateString(msg)
	reasonMem := pdk.AllocateString(reason)
	var outcomeVal uint64
	if outcome {
		outcomeVal = 1
	}
	assert(msgMem.Offset(), outcomeVal, reasonMem.Offset())
	msgMem.Free()
	reasonMem.Free()
}

// Assert that `x` and `y` are equal, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func AssertEq(msg string, x any, y any) {
	Assert(msg, x == y, fmt.Sprintf("Expected %v == %v", x, y))
}

// Assert that `x` and `y` are not equal, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func AssertNe(msg string, x any, y any) {
	Assert(msg, x != y, fmt.Sprintf("Expected %v != %v", x, y))
}

// Assert that `x` is greater than `y`, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func AssertGt[T cmp.Ordered](msg string, x T, y T) {
	Assert(msg, x > y, fmt.Sprintf("Expected %v > %v", x, y))
}

// Assert that `x` is greater than or equal to `y`, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func AssertGte[T cmp.Ordered](msg string, x T, y T) {
	Assert(msg, x >= y, fmt.Sprintf("Expected %v >= %v", x, y))
}

// Assert that `x` is less than `y`, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func AssertLt[T cmp.Ordered](msg string, x T, y T) {
	Assert(msg, x < y, fmt.Sprintf("Expected %v < %v", x, y))
}

// Assert that `x` is less than or equal to `y`, naming the assertion with `msg`, which will be used as a label in the CLI runner.
func AssertLte[T cmp.Ordered](msg string, x T, y T) {
	Assert(msg, x <= y, fmt.Sprintf("Expected %v <= %v", x, y))
}

// Reset the loaded plugin, clearing all state.
func Reset() {
	reset()
}

// Create a new test group. NOTE: these cannot be nested and starting a new group will end the last one.
func startGroup(name string) {
	groupNameMem := pdk.AllocateString(name)
	group(groupNameMem.Offset())
	groupNameMem.Free()
}

// Run a test group, resetting the plugin before and after the group is run.
func Group(name string, tests func()) {
	Reset()
	startGroup(name)
	tests()
	Reset()
}

// Read the mock test input provided by the test runner, returns a Memory object.
// This input is defined in an xtp.toml file, or by the --mock-input-data or --mock-input-file flags.
func MockInput() pdk.Memory {
	outputPtr := mock_input()
	return pdk.FindMemory(outputPtr)
}

// Read the mock test input provided by the test runner, returns the input as []byte.
// This input is defined in an xtp.toml file, or by the --mock-input-data or --mock-input-file flags.
func MockInputBytes() []byte {
	inputMem := MockInput()
	defer inputMem.Free()
	return inputMem.ReadBytes()
}

// Read the mock test input provided by the test runner, returns the input as a string.
// This input is defined in an xtp.toml file, or by the --mock-input-data or --mock-input-file flags.
func MockInputString() string {
	return string(MockInputBytes())
}
