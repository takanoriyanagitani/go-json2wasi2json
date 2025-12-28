package wazero

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/takanoriyanagitani/go-json2wasi2json"
)

func mustGetWasm(t testing.TB) []byte {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get working directory: %v", err)
	}
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(wd)))
	wasmPath := filepath.Join(projectRoot, "ident.wasm")

	wasm, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("Failed to read ident.wasm from %s: %v", wasmPath, err)
	}
	return wasm
}

func TestWazeroImplementation(t *testing.T) {
	wasm := mustGetWasm(t)
	ctx := context.Background()

	t.Run("raw func roundtrip", func(t *testing.T) {
		t.Parallel()
		fn, close, err := New(ctx, wasm, json2wasi2json.Limits{})
		if err != nil {
			t.Fatalf("failed to create new wazero function: %v", err)
		}
		defer func() {
			if err := close(); err != nil {
				t.Errorf("wazero closer failed: %v", err)
			}
		}()

		input := []byte(`{"test": "wazero"}`)
		expected := []byte(`{"test": "wazero"}`)

		output, err := fn(ctx, input)
		if err != nil {
			t.Fatalf("RawFunc returned an unexpected error: %v", err)
		}

		if !bytes.Equal(output, expected) {
			t.Errorf("wazero RawFunc output mismatch: got %s, want %s", output, expected)
		}
	})
}

func BenchmarkIdentRawSafe(b *testing.B) {
	wasm := mustGetWasm(b)
	ctx := context.Background()
	fn, close, err := New(ctx, wasm, json2wasi2json.Limits{})
	if err != nil {
		b.Fatalf("failed to create new wazero function: %v", err)
	}
	defer func() {
		if err := close(); err != nil {
			b.Errorf("wazero closer failed: %v", err)
		}
	}()

	input := []byte(`{"message":"benchmark"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fn(ctx, input)
		if err != nil {
			b.Fatalf("Function returned an error during benchmark: %v", err)
		}
	}
}

type BenchmarkPayload struct {
	Message string `json:"message"`
}

func BenchmarkIdentStructSafe(b *testing.B) {
	wasm := mustGetWasm(b)
	ctx := context.Background()
	rawFunc, close, err := New(ctx, wasm, json2wasi2json.Limits{})
	if err != nil {
		b.Fatalf("failed to create new wazero function: %v", err)
	}
	defer func() {
		if err := close(); err != nil {
			b.Errorf("wazero closer failed: %v", err)
		}
	}()

	structFunc := json2wasi2json.ToPureFunc[BenchmarkPayload, BenchmarkPayload](rawFunc)
	input := BenchmarkPayload{Message: "benchmark"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := structFunc(ctx, input)
		if err != nil {
			b.Fatalf("Function returned an error during benchmark: %v", err)
		}
	}
}

func BenchmarkIdentRawUnsafe(b *testing.B) {
	wasm := mustGetWasm(b)
	ctx := context.Background()
	limits := json2wasi2json.Limits{DisableContextDoneChecks: true}
	fn, close, err := New(ctx, wasm, limits)
	if err != nil {
		b.Fatalf("failed to create new wazero function: %v", err)
	}
	defer func() {
		if err := close(); err != nil {
			b.Errorf("wazero closer failed: %v", err)
		}
	}()

	input := []byte(`{"message":"benchmark"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fn(ctx, input)
		if err != nil {
			b.Fatalf("Function returned an error during benchmark: %v", err)
		}
	}
}

func BenchmarkIdentStructUnsafe(b *testing.B) {
	wasm := mustGetWasm(b)
	ctx := context.Background()
	limits := json2wasi2json.Limits{DisableContextDoneChecks: true}
	rawFunc, close, err := New(ctx, wasm, limits)
	if err != nil {
		b.Fatalf("failed to create new wazero function: %v", err)
	}
	defer func() {
		if err := close(); err != nil {
			b.Errorf("wazero closer failed: %v", err)
		}
	}()

	structFunc := json2wasi2json.ToPureFunc[BenchmarkPayload, BenchmarkPayload](rawFunc)
	input := BenchmarkPayload{Message: "benchmark"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := structFunc(ctx, input)
		if err != nil {
			b.Fatalf("Function returned an error during benchmark: %v", err)
		}
	}
}
