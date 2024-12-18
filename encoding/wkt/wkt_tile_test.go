package wkt

import (
	"bytes"
	"testing"

	gtesting "github.com/hahaking119/geom/testing"
)

func init() {
	gtesting.CompileTiles(DecodeString)
}

func BenchmarkEncodeTile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncodeBytes(gtesting.Tiles()[0])
	}
}

func BenchmarkEncodeTilePrealloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// the encoded wkt is ~32MB
		buf := bytes.NewBuffer(make([]byte, 0, (1<<20)*32))
		enc := NewDefaultEncoder(buf)
		enc.Encode(gtesting.Tiles()[0])
	}
}

func BenchmarkEncodeTileNoprealloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(make([]byte, 0, 0))
		enc := NewDefaultEncoder(buf)
		enc.Encode(gtesting.Tiles()[0])
	}
}
