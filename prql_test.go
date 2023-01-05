package prql

import (
	"context"
	"os"
	"testing"
)

func TestWasi(t *testing.T) {
	ctx := context.Background()
	w, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close(ctx)

	testCases := map[string]struct {
		input     string
		shouldErr bool
	}{
		"empty string": {
			input:     "",
			shouldErr: true,
		},
		"invalid prql": {
			input:     "foo",
			shouldErr: true,
		},
		"valid prql": {
			input:     `from employees`,
			shouldErr: false,
		},
		"example query": {
			input:     MustFromFile("testdata/input.prql"),
			shouldErr: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			_, err := w.Compile(ctx, tc.input)
			if tc.shouldErr && err == nil {
				t.Fatalf("expected err, got nil")
			} else if !tc.shouldErr && err != nil {
				t.Fatalf("expected nil, got err: %v", err)
			}
		})
	}
}

func MustFromFile(name string) string {
	content, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(content)
}
