package roaring

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestBitmap_Reset verifies that the Reset method always restores a Bitmap
// to a state indistinguishable from a new one.
func TestBitmap_Reset(t *testing.T) {
	bm := NewBitmapWithDefaultPooling(10)

	i := uint64(0)
	for {
		if i >= 10000 {
			break
		}

		bm.Add(10000 * i)
		i++
	}
	untouched := NewBitmapWithDefaultPooling(10)
	bm.Reset()

	if !cmp.Equal(untouched, bm, cmp.Comparer(func(x, y Bitmap) bool {
		if x.Containers.Size() != y.Containers.Size() {
			return false
		}
		if x.Containers.Count() != y.Containers.Count() {
			return false
		}
		if x.opN != y.opN {
			return false
		}
		if x.OpWriter != y.OpWriter {
			return false
		}
		return true
	})) {
		t.Fatalf("Reset bitmap: %+v is not identical to new bitmap: %+v", bm, untouched)
	}
}
