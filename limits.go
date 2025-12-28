package json2wasi2json

// Limits defines constraints for a WASM execution.
type Limits struct {
	// MemoryMaxPages specifies the maximum number of memory pages (64KB each)
	// the module can use. If zero, a runtime-specific default is used.
	MemoryMaxPages uint32
	// DisableContextDoneChecks disables the WithCloseOnContextDone feature.
	// This improves performance but is less safe for untrusted modules.
	// It is false by default, making the safe behavior the default.
	DisableContextDoneChecks bool
}

const WasmPageSizeKiB = 64
const KiBytesInMiByte = 1024
const WasmPagesInMiB = KiBytesInMiByte / WasmPageSizeKiB

// MemoryMaxMiB converts a memory limit in MiB to pages.
func (l Limits) MemoryMaxMiB(mib uint32) Limits {
	l.MemoryMaxPages = mib * WasmPagesInMiB
	return l
}
