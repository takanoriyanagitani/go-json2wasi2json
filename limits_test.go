package json2wasi2json

import "testing"

func TestLimitsMemoryMaxMiB(t *testing.T) {
	t.Parallel()
	t.Run("zero mib", func(t *testing.T) {
		t.Parallel()
		l := Limits{}.MemoryMaxMiB(0)
		if l.MemoryMaxPages != 0 {
			t.Errorf("Expected 0 pages for 0 MiB, got %d", l.MemoryMaxPages)
		}
	})

	t.Run("one mib", func(t *testing.T) {
		t.Parallel()
		l := Limits{}.MemoryMaxMiB(1)
		expected := uint32(1 * WasmPagesInMiB)
		if l.MemoryMaxPages != expected {
			t.Errorf("Expected %d pages for 1 MiB, got %d", expected, l.MemoryMaxPages)
		}
	})

	t.Run("64 mib", func(t *testing.T) {
		t.Parallel()
		l := Limits{}.MemoryMaxMiB(64)
		expected := uint32(64 * WasmPagesInMiB)
		if l.MemoryMaxPages != expected {
			t.Errorf("Expected %d pages for 64 MiB, got %d", expected, l.MemoryMaxPages)
		}
	})
}
