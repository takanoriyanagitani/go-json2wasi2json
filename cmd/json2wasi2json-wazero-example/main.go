package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/takanoriyanagitani/go-json2wasi2json"
	"github.com/takanoriyanagitani/go-json2wasi2json/wasi/runtime/wazero"
)

func must(e error) {
	if nil != e {
		log.Fatalln(e)
	}
}

// MessageInput struct for rs-msgcount.wasm
type MessageInput struct {
	Message string `json:"message"`
}

// MessageOutput struct for rs-msgcount.wasm
type MessageOutput struct {
	MessageLength uint32 `json:"message_length"`
}

func main() {
	// This example uses rs-msgcount.wasm, which expects a JSON object with a "message" field
	// and returns a JSON object with a "message_length" field.
	wasm, e := os.ReadFile("cmd/json2wasi2json-wazero-example/rs-msgcount/rs-msgcount.wasm")
	must(e)

	// --- Configure Limits (Memory Limit) ---
	// Set a memory limit of 64 MiB
	limits := json2wasi2json.Limits{}.MemoryMaxMiB(64)
	fmt.Printf("Configured WASM Memory Limit: %d pages (%d MiB)\n", limits.MemoryMaxPages, limits.MemoryMaxPages/json2wasi2json.WasmPagesInMiB)

	ctx := context.Background()
	fn, close, e := wazero.New(ctx, wasm, limits) // Pass limits to New
	must(e)
	defer func() {
		e := close()
		if nil != e {
			log.Printf("close failed: %v\n", e)
		}
	}()

	// --- Demonstrate Raw byte processing with rs-msgcount's expected input ---
	rawInput := []byte(`{"message": "hi"}`)
	fmt.Printf("Raw Input:  %s\n", rawInput)

	rawOutput, e := fn(ctx, rawInput)
	must(e)

	fmt.Printf("Raw Output: %s\n", rawOutput)

	// --- Demonstrate JsonMap processing with rs-msgcount ---
	jsonMapFunc := json2wasi2json.RawFunc(fn).ToPureFuncJsonMap(
		json2wasi2json.Bytes2JsonMap,
		json2wasi2json.JsonMap2Bytes,
	)

	jsonMapInput := json2wasi2json.JsonMap{"message": "msg"}
	fmt.Printf("JsonMap Input: %#v\n", jsonMapInput)

	jsonMapOutput, e := jsonMapFunc(ctx, jsonMapInput)
	must(e)

	fmt.Printf("JsonMap Output: %#v\n", jsonMapOutput)

	// --- Demonstrate Struct-based processing with rs-msgcount ---
	structFunc := json2wasi2json.ToPureFunc[MessageInput, MessageOutput](fn)

	structInput := MessageInput{Message: "abc"}
	fmt.Printf("Struct Input: %#v\n", structInput)

	structOutput, e := structFunc(ctx, structInput)
	must(e)

	fmt.Printf("Struct Output: %#v\n", structOutput)

	// --- Demonstrate Timeout ---
	fmt.Println("\n--- Demonstrating Timeout ---")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Use the RawFunc for timeout demonstration
	_, timeoutErr := fn(timeoutCtx, rawInput)
	if timeoutErr != nil {
		if errors.Is(timeoutErr, context.DeadlineExceeded) {
			fmt.Printf("Function call timed out as expected: %v\n", timeoutErr)
		} else {
			fmt.Printf("Function call failed with unexpected error during timeout test: %v\n", timeoutErr)
		}
	} else {
		fmt.Println("Function completed within timeout, which was not expected.")
	}
}
