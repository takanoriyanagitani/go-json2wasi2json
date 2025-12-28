package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/takanoriyanagitani/go-json2wasi2json"
	"github.com/takanoriyanagitani/go-json2wasi2json/wasi/runtime/wazero"
)

func must(e error) {
	if nil != e {
		log.Fatalln(e)
	}
}

// IdentPayload is a simple struct for demonstrating struct-based input/output.
type IdentPayload struct {
	Message string `json:"message"`
	Value   int    `json:"value"`
}

func main() {
	// This example assumes a WASM file named `ident.wasm` exists in the current directory.
	// The `ident.wasm` module should be a WASI command that reads from stdin
	// and writes the exact same content to stdout.
	wasm, e := os.ReadFile("ident.wasm")
	must(e)

	ctx := context.Background()
	fn, close, e := wazero.New(ctx, wasm, json2wasi2json.Limits{})
	must(e)
	defer func() {
		e := close()
		if nil != e {
			log.Printf("close failed: %v\n", e)
		}
	}()

	var input []byte = []byte(`{"hello": "world"}`)
	fmt.Printf("Raw Input:  %s\n", input)

	output, e := fn(ctx, input)
	must(e)

	fmt.Printf("Raw Output: %s\n", output)

	// --- Demonstrate JsonMap processing ---
	jsonMapFunc := json2wasi2json.RawFunc(fn).ToPureFuncJsonMap(
		json2wasi2json.Bytes2JsonMap,
		json2wasi2json.JsonMap2Bytes,
	)

	jsonMapInput := json2wasi2json.JsonMap{"message": "hello from JsonMap"}
	fmt.Printf("JsonMap Input: %#v\n", jsonMapInput)

	jsonMapOutput, e := jsonMapFunc(ctx, jsonMapInput)
	must(e)

	fmt.Printf("JsonMap Output: %#v\n", jsonMapOutput)

	// --- Demonstrate Struct-based processing ---
	structFunc := json2wasi2json.ToPureFunc[IdentPayload, IdentPayload](fn)

	structInput := IdentPayload{Message: "hello from struct", Value: 123}
	fmt.Printf("Struct Input: %#v\n", structInput)

	structOutput, e := structFunc(ctx, structInput)
	must(e)

	fmt.Printf("Struct Output: %#v\n", structOutput)
}
