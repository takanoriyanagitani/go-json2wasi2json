package wazero

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/takanoriyanagitani/go-json2wasi2json"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

type UUID struct{ uuid.UUID }

func (u UUID) String() string { return u.UUID.String() }

func (u UUID) ToInstanceName() string { return "instance-" + u.String() }

func UUID7() (UUID, error) {
	uid, e := uuid.NewV7()
	if e != nil {
		return UUID{}, fmt.Errorf("failed to create V7 UUID: %w", e)
	}
	return UUID{UUID: uid}, nil
}

type WasmConfig struct{ wazero.ModuleConfig }

func (c WasmConfig) WithName(name string) WasmConfig {
	return WasmConfig{
		ModuleConfig: c.ModuleConfig.WithName(name),
	}
}

func (c WasmConfig) WithID(id UUID) WasmConfig {
	name := id.ToInstanceName()
	return c.WithName(name)
}

func (c WasmConfig) WithAutoID() (WasmConfig, error) {
	id, e := UUID7()
	if e != nil {
		return WasmConfig{}, fmt.Errorf("failed to get auto ID: %w", e)
	}
	neo := c.WithID(id)
	return neo, nil
}

func (c WasmConfig) WithReader(rdr io.Reader) WasmConfig {
	return WasmConfig{
		ModuleConfig: c.ModuleConfig.WithStdin(rdr),
	}
}

func (c WasmConfig) WithWriter(wtr io.Writer) WasmConfig {
	return WasmConfig{
		ModuleConfig: c.ModuleConfig.WithStdout(wtr),
	}
}

var ErrInstantiate error = errors.New("unable to instantiate wasm module")

func newFunc(
	r wazero.Runtime,
	compiled wazero.CompiledModule,
	conf wazero.ModuleConfig,
) json2wasi2json.RawFunc {
	return func(ctx context.Context, input []byte) ([]byte, error) {
		var stdout bytes.Buffer
		config := conf.
			WithStdin(bytes.NewReader(input)).
			WithStdout(&stdout)
		instance, e := r.InstantiateModule(ctx, compiled, config)
		if nil != e {
			return nil, fmt.Errorf("%w: %w", ErrInstantiate, e)
		}

		e = instance.Close(ctx)
		var exitErr *sys.ExitError
		if errors.As(e, &exitErr) {
			if 0 != exitErr.ExitCode() {
				return stdout.Bytes(), exitErr
			}
			return stdout.Bytes(), nil
		}
		return stdout.Bytes(), e
	}
}

func New(
	ctx context.Context,
	wasmBytes []byte,
	limits json2wasi2json.Limits,
) (json2wasi2json.RawFunc, func() error, error) {
	runtimeConfig := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(!limits.DisableContextDoneChecks)

	if limits.MemoryMaxPages > 0 {
		runtimeConfig = runtimeConfig.WithMemoryLimitPages(limits.MemoryMaxPages)
	}

	rt := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)
	closer := func() error { return rt.Close(ctx) }

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	compiled, e := rt.CompileModule(ctx, wasmBytes)
	if nil != e {
		return nil, closer, e
	}

	conf := wazero.NewModuleConfig()

	return newFunc(rt, compiled, conf), closer, nil
}
