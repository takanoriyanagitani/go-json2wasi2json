#!/bin/sh

set -eu

cd "$(dirname "$0")"

PKG=./wasi/runtime/wazero

bench_raw_safe() {
  go test -bench='^BenchmarkIdentRawSafe$' "$PKG"
}

bench_struct_safe() {
  go test -bench='^BenchmarkIdentStructSafe$' "$PKG"
}

bench_raw_unsafe() {
  go test -bench='^BenchmarkIdentRawUnsafe$' "$PKG"
}

bench_struct_unsafe() {
  go test -bench='^BenchmarkIdentStructUnsafe$' "$PKG"
}

all() {
  go test -bench=. "$PKG"
}

main() {
  if [ $# -eq 0 ]; then
    all
    return
  fi

  "$@"
}

main "$@"