package main

import (
	pdk "github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user kv_read
func kv_read(key uint64) uint64

//go:wasmimport extism:host/user kv_write
func kv_write(key uint64, value uint64)

//go:export run
func run() int32 {
	key := pdk.AllocateString("key")
	value := pdk.AllocateString("value")
	kv_write(key.Offset(), value.Offset())

	readVal := kv_read(key.Offset())
	if readVal != 0 {
		readValMem := pdk.FindMemory(readVal)
		varVal := pdk.GetVar("key")
		pdk.SetVar("key", append(varVal, readValMem.ReadBytes()...))
	} else {
		pdk.SetVar("key", []byte(""))
	}

	pdk.Output(pdk.GetVar("key"))
	return 0
}

func main() {}
