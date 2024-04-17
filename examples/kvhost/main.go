package main

import (
	pdk "github.com/extism/go-pdk"
)

var kv map[string]string = make(map[string]string)

//go:export kv_read
func kv_read(key uint64) uint64 {
	keyMem := pdk.FindMemory(key)
	k := string(keyMem.ReadBytes())
	v, ok := kv[k]
	if !ok {
		return 0
	}

	valMem := pdk.AllocateString(v)
	return valMem.Offset()
}

//go:export kv_write
func kv_write(key uint64, value uint64) {
	keyMem := pdk.FindMemory(key)
	valueMem := pdk.FindMemory(value)
	k := string(keyMem.ReadBytes())
	v := string(valueMem.ReadBytes())
	kv[k] = v
}

func main() {}
