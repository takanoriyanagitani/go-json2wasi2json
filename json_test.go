package json2wasi2json

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBytes2JsonMap(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		t.Parallel()
		input := []byte(`{"hello": "world", "value": 99}`)
		expected := JsonMap{"hello": "world", "value": 99.0} // JSON numbers are float64 by default

		result, err := Bytes2JsonMap(input)
		if err != nil {
			t.Fatalf("Bytes2JsonMap returned an unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Bytes2JsonMap output mismatch: got %#v, want %#v", result, expected)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		input := []byte(`not valid json`)
		_, err := Bytes2JsonMap(input)
		if err == nil {
			t.Fatal("Bytes2JsonMap did not return an error for invalid JSON")
		}
	})
}

func TestJsonMap2Bytes(t *testing.T) {
	t.Parallel()
	input := JsonMap{"message": "test", "value": 123.0}

	result, err := JsonMap2Bytes(input)
	if err != nil {
		t.Fatalf("JsonMap2Bytes returned an unexpected error: %v", err)
	}

	// Note: JSON marshaling doesn't guarantee key order, so a direct string comparison is fragile.
	// A better test is to unmarshal the result and check for equality.
	var roundtripMap JsonMap
	if err := json.Unmarshal(result, &roundtripMap); err != nil {
		t.Fatalf("Failed to unmarshal result for comparison: %v", err)
	}
	if !reflect.DeepEqual(input, roundtripMap) {
		t.Errorf("JsonMap2Bytes roundtrip failed. Expected %s to produce %#v", string(result), input)
	}
}
