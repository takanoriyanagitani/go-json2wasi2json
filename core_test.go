package json2wasi2json

import (
	"context"
	"reflect"
	"testing"
)

// mockIdentRawFunc is a mock RawFunc that simulates an identity WASM call.
// It simply returns the input bytes, mimicking a WASM module that echoes its stdin.
func mockIdentRawFunc(ctx context.Context, input []byte) ([]byte, error) {
	return input, nil
}

// TestPayload is a simple struct used for testing the generic adapters.
type TestPayload struct {
	Message string `json:"message"`
	Value   int    `json:"value"`
}

func TestToPureFunc(t *testing.T) {
	t.Parallel()
	// 1. Create the high-level struct-based function from the mock RawFunc.
	structFunc := ToPureFunc[TestPayload, TestPayload](mockIdentRawFunc)

	// 2. Define the input struct.
	input := TestPayload{Message: "hello", Value: 42}
	expected := TestPayload{Message: "hello", Value: 42}

	// 3. Call the function.
	output, err := structFunc(context.Background(), input)
	if err != nil {
		t.Fatalf("ToPureFunc returned an unexpected error: %v", err)
	}

	// 4. Assert the output is correct.
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("ToPureFunc output mismatch: got %#v, want %#v", output, expected)
	}
}

func TestToPureFuncJsonMap(t *testing.T) {
	t.Parallel()
	// 1. Create the high-level map-based function.
	jsonMapFunc := RawFunc(mockIdentRawFunc).ToPureFuncJsonMap(
		Bytes2JsonMap,
		JsonMap2Bytes,
	)

	// 2. Define the input map.
	input := JsonMap{"message": "world", "value": 123.0} // Use float64 for numbers, as is standard for json.Unmarshal
	expected := JsonMap{"message": "world", "value": 123.0}

	// 3. Call the function.
	output, err := jsonMapFunc(context.Background(), input)
	if err != nil {
		t.Fatalf("ToPureFuncJsonMap returned an unexpected error: %v", err)
	}

	// 4. Assert the output is correct.
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("ToPureFuncJsonMap output mismatch: got %#v, want %#v", output, expected)
	}
}
