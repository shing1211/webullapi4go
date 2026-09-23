package data_test

import (
	"encoding/json"
	"testing"

	"github.com/shing1211/webullapi4go/data"
)

func FuzzDecodeQuote(f *testing.F) {
	f.Add(`{"symbol":"AAPL"}`)
	f.Fuzz(func(t *testing.T, input string) {
		var v data.Quote
		_ = json.Unmarshal([]byte(input), &v)
	})
}

func FuzzDecodeSnapshot(f *testing.F) {
	f.Add(`{"symbol":"AAPL"}`)
	f.Fuzz(func(t *testing.T, input string) {
		var v data.Snapshot
		_ = json.Unmarshal([]byte(input), &v)
	})
}

func FuzzDecodeTick(f *testing.F) {
	f.Add(`{"symbol":"AAPL"}`)
	f.Fuzz(func(t *testing.T, input string) {
		var v data.StockTicks
		_ = json.Unmarshal([]byte(input), &v)
	})
}
