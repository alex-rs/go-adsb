// Copyright 2024 Collin Kreklow
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package beast_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/alex-rs/go-adsb/beast"
	"github.com/alex-rs/go-adsb/internal/benchdata"
)

func BenchmarkDecoderStream(b *testing.B) {
	stream := benchdata.BeastStream(16) // 3 frame types * 16 = 48 frames per run

	b.ReportAllocs()
	b.SetBytes(int64(len(stream)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(stream)
		d := beast.NewDecoder(r)
		d.StripEscape = true
		d.ResyncOnError = true

		f := new(beast.Frame)

		for {
			if err := d.Decode(f); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				b.Fatalf("unexpected decode error: %v", err)
			}
		}
	}
}

func BenchmarkDecoderNoisyStream(b *testing.B) {
	stream := benchdata.NoisyBeastStream(32)

	b.ReportAllocs()
	b.SetBytes(int64(len(stream)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(stream)
		d := beast.NewDecoder(r)
		d.StripEscape = true

		f := new(beast.Frame)

		for {
			if err := d.Decode(f); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				// stop after first decode error; stream intentionally
				// contains invalid data.
				break
			}
		}
	}
}
