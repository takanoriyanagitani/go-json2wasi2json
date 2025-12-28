package json2wasi2json

import (
	"context"
	"encoding/json"
)

type PureFuncJson[I, O any] func(context.Context, I) (O, error)

type JsonMap map[string]any

type PureFuncJsonMap PureFuncJson[JsonMap, JsonMap]

type RawFunc func(context.Context, []byte) ([]byte, error)

type ToMap func([]byte) (JsonMap, error)
type ToBytes func(JsonMap) ([]byte, error)

func (r RawFunc) ToPureFuncJsonMap(b2j ToMap, j2b ToBytes) PureFuncJsonMap {
	return func(ctx context.Context, input JsonMap) (JsonMap, error) {
		jbytes, e := j2b(input)
		if nil != e {
			return nil, e
		}

		raw, e := r(ctx, jbytes)
		if nil != e {
			return nil, e
		}

		return b2j(raw)
	}
}

// ToPureFunc converts a low-level RawFunc into a high-level generic function
// that works directly with JSON-compatible structs or types.
// It handles the JSON marshalling of the input and unmarshalling of the output.
func ToPureFunc[I, O any](r RawFunc) PureFuncJson[I, O] {
	return func(ctx context.Context, input I) (O, error) {
		var output O // The zero value of the output struct

		// Marshal the input struct to JSON bytes
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return output, err
		}

		// Call the low-level RawFunc
		outputBytes, err := r(ctx, inputBytes)
		if err != nil {
			return output, err
		}

		// Unmarshal the output bytes into the output struct
		err = json.Unmarshal(outputBytes, &output)
		if err != nil {
			return output, err
		}

		return output, nil
	}
}
