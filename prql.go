package prql

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"hash"
	"strings"

	"github.com/segmentio/fasthash/fnv1a"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Engine interface {
	Compile(context.Context, string) (string, error)
}

var _ Engine = (*WasiEngine)(nil)

//go:embed testdata/prql-wasi.wasm
var prqlWasi []byte

type WasiEngine struct {
	code wazero.CompiledModule
	r    wazero.Runtime
	h    hash.Hash
}

// Compile compiles a prql query to a sql query
func (e *WasiEngine) Compile(ctx context.Context, query string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("prql query must not be empty")
	}

	// hash query for concurrent a
	h1 := fnv1a.HashString64(query)
	name := fmt.Sprintf("id-%x", h1)

	in := strings.NewReader(query)
	out := new(bytes.Buffer)
	// we know our wasi program doesn't write to stderr
	// so we skip configuring it

	config := wazero.NewModuleConfig().
		WithStdout(out).
		WithStdin(in).
		WithName(name)

	mod, err := e.r.InstantiateModule(ctx, e.code, config)
	if err != nil {
		return "", err
	}

	mod.Close(ctx)
	return strings.TrimSpace(out.String()), nil
}

// Close closes the underlying wazero runtime
func (e *WasiEngine) Close(ctx context.Context) error {
	return e.r.Close(ctx)
}

// New instantiates a new wasi runtime and precompiles the embedded wasi file
func New(ctx context.Context) (*WasiEngine, error) {
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	code, err := r.CompileModule(ctx, prqlWasi)
	if err != nil {
		return nil, err
	}

	return &WasiEngine{
		r:    r,
		code: code,
		h:    sha256.New(),
	}, nil
}
