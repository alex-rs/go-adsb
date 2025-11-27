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

// Package benchdata exposes shared fixtures for benchmarks.
package benchdata

import (
	"bytes"
	"encoding/hex"
)

// MustHex converts a hex string to bytes and panics on error. Intended
// for test and benchmark fixtures only.
func MustHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}

// BeastFrames contains valid Beast messages that exercise common frame
// types. Data is in wire format (including escape bytes).
var BeastFrames = [][]byte{
	// Type 1 Mode A/C
	MustHex("1a311a1af933baf325c45047"),
	// Type 2 ADS-B / Mode S
	MustHex("1a321a1af933baf325c45da99adad95ff6"),
	// Type 3 ADS-B / Mode S with escaped 0x1a bytes
	MustHex("1a331a1af933bbc63ec68f1a1a9ada58b98446e703357e241a1a"),
}

var noiseChunks = [][]byte{
	MustHex("00ff1337"),
	[]byte("garbage"),
	MustHex("deadbeef1a99ff"),
}

// ADSBMessages contains representative ADS-B payloads across downlink
// formats used in benchmarks.
var ADSBMessages = [][]byte{
	// DF17, identity
	MustHex("8dacf84e23101332cf3ca037ef13"),
	// DF17, global position message 1
	MustHex("8da8028758ab0028de078689d437"),
	// DF17, global position message 2
	MustHex("8da8028758ab07b0b8876e81eb25"),
	// DF20, Comm-B altitude reply with callsign
	MustHex("a0000f9820057273df8d20e2cf30"),
	// DF5, identity reply
	MustHex("28001b0601970d"),
}

// BeastStream returns a concatenated Beast stream made up of the
// packaged frames, repeated n times to build longer inputs.
func BeastStream(n int) []byte {
	var buf bytes.Buffer

	for i := 0; i < n; i++ {
		for _, f := range BeastFrames {
			buf.Write(f)
		}
	}

	return buf.Bytes()
}

// NoisyBeastStream returns a Beast stream interspersed with invalid
// bytes to exercise decoder resync logic. The pattern is deterministic
// for repeatability.
func NoisyBeastStream(n int) []byte {
	var buf bytes.Buffer

	for i := 0; i < n; i++ {
		buf.Write(noiseChunks[i%len(noiseChunks)])
		if i%3 == 0 {
			// insert a truncated header to force resync
			buf.Write([]byte{0x1a, 0x33, 0xff})
		}
		buf.Write(BeastFrames[i%len(BeastFrames)])
		buf.Write(noiseChunks[(i+1)%len(noiseChunks)])
	}

	return buf.Bytes()
}
