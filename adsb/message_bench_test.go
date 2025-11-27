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

package adsb_test

import (
	"fmt"
	"testing"

	"github.com/alex-rs/go-adsb/adsb"
	"github.com/alex-rs/go-adsb/internal/benchdata"
)

func benchmarkMessage(b *testing.B, payload []byte, fn func(m *adsb.Message) error) {
	b.Helper()

	msg := new(adsb.Message)

	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := msg.UnmarshalBinary(payload); err != nil {
			b.Fatalf("unmarshal failed: %v", err)
		}

		if err := fn(msg); err != nil {
			b.Fatalf("operation failed: %v", err)
		}
	}
}

func BenchmarkMessageICAO(b *testing.B) {
	payload := benchdata.ADSBMessages[0] // DF17 identity

	benchmarkMessage(b, payload, func(m *adsb.Message) error {
		_, err := m.ICAO()
		if err != nil {
			return fmt.Errorf("icao: %w", err)
		}

		return nil
	})
}

func BenchmarkMessageAltitude(b *testing.B) {
	payload := benchdata.ADSBMessages[1] // DF17 global position with altitude

	benchmarkMessage(b, payload, func(m *adsb.Message) error {
		_, err := m.Alt()
		if err != nil {
			return fmt.Errorf("altitude: %w", err)
		}

		return nil
	})
}

func BenchmarkMessageCallsign(b *testing.B) {
	payload := benchdata.ADSBMessages[3] // DF20 with callsign
	callBuf := make([]byte, 0, 8)

	benchmarkMessage(b, payload, func(m *adsb.Message) error {
		_, err := m.CallBytes(callBuf[:0])
		if err != nil {
			return fmt.Errorf("callsign: %w", err)
		}

		return nil
	})
}

func BenchmarkMessageGlobalPosition(b *testing.B) {
	// Pair of DF17 global position messages with alternating format bits.
	payloadEven := benchdata.ADSBMessages[1]
	payloadOdd := benchdata.ADSBMessages[2]

	var cprEven, cprOdd adsb.CPR
	dst := make([]float64, 2)

	m1 := new(adsb.Message)
	m2 := new(adsb.Message)

	b.ReportAllocs()
	b.SetBytes(int64(len(payloadEven) + len(payloadOdd)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := m1.UnmarshalBinary(payloadEven); err != nil {
			b.Fatalf("unmarshal m1 failed: %v", err)
		}

		c1, err := m1.CPRInto(&cprEven)
		if err != nil {
			b.Fatalf("cpr m1 failed: %v", err)
		}

		if err := m2.UnmarshalBinary(payloadOdd); err != nil {
			b.Fatalf("unmarshal m2 failed: %v", err)
		}

		c2, err := m2.CPRInto(&cprOdd)
		if err != nil {
			b.Fatalf("cpr m2 failed: %v", err)
		}

		if _, err := adsb.DecodeGlobalPositionInto(c1, c2, dst); err != nil {
			b.Fatalf("decode global position failed: %v", err)
		}
	}
}
