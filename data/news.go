// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// pathNewsSummaries is the news-summary endpoint. Unlike the other market-data
// endpoints it replies with a Server-Sent Events (SSE) stream.
//
// Reference: https://developer.webull.hk/apis/docs/reference/news-summary.md
const pathNewsSummaries = "/market-data/news/summaries/get"

// NewsCategorySymbols groups the security symbols of one category in a news
// summary request.
type NewsCategorySymbols struct {
	// Category is the security type. Only [StockCategoryUS] is currently
	// accepted.
	Category StockCategory `json:"category"`
	// Symbols lists the security symbols, for example ["AAPL", "GOOG"].
	Symbols []string `json:"symbols"`
}

// NewsSummaryParam parameterizes [Client.GetNewsSummary]. CategorySymbols is
// required.
type NewsSummaryParam struct {
	// CategorySymbols lists the security symbols to summarize, grouped by
	// category. Required.
	CategorySymbols []NewsCategorySymbols `json:"category_symbols"`
	// Lang is the response language. Only "en" is currently supported; empty
	// means the server default.
	Lang string `json:"lang,omitempty"`
}

// NewsSummaryText is a single cell of a news-summary table event.
type NewsSummaryText struct {
	// Text is the cell text, which may contain Markdown.
	Text string `json:"text"`
}

// NewsSummaryEvent is one event in the news-summary SSE stream.
//
// The stream emits a single "meta" event that carries the conversation
// identifiers in RawArgs, followed by any number of "text" events (Markdown in
// Message) and "table" events (Headers and Rows). Unknown event types are
// returned with only Type populated so callers can ignore or log them.
type NewsSummaryEvent struct {
	// Type is the event kind: "meta", "text", or "table".
	Type string `json:"type"`
	// Message is the Markdown text of a "text" event.
	Message string `json:"message,omitempty"`
	// Args is the raw JSON metadata of a "meta" event, for example
	// {"sessionId":"1","convId":451107711450219}. It is left undecoded so
	// callers can read the fields they need.
	Args json.RawMessage `json:"args,omitempty"`
	// Headers is the column headers of a "table" event.
	Headers []NewsSummaryText `json:"headers,omitempty"`
	// Rows is the row data of a "table" event.
	Rows [][]NewsSummaryText `json:"rows,omitempty"`
}

// NewsSummaryStream reads a news-summary SSE response. Call [Next] until it
// returns [io.EOF], then [Close]. A stream is NOT safe for concurrent use.
type NewsSummaryStream struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
}

// GetNewsSummary opens a news-summary stream for the requested symbols. The
// caller must close the returned stream. A non-2xx response is returned as a
// typed error instead of a stream.
//
// Reference: https://developer.webull.hk/apis/docs/reference/news-summary.md
func (c *Client) GetNewsSummary(ctx context.Context, params NewsSummaryParam) (*NewsSummaryStream, error) {
	resp, err := c.core.DoStream(ctx, http.MethodPost, pathNewsSummaries, params)
	if err != nil {
		return nil, err
	}
	return newNewsSummaryStream(resp.Body), nil
}

// newNewsSummaryStream wraps body in an SSE line scanner.
func newNewsSummaryStream(body io.ReadCloser) *NewsSummaryStream {
	scanner := bufio.NewScanner(body)
	// News summaries can carry long Markdown lines; raise the default 64 KiB
	// token limit to 1 MiB.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return &NewsSummaryStream{body: body, scanner: scanner}
}

// Next returns the next news-summary event, or [io.EOF] when the stream ends.
// SSE comment lines and field lines other than "data" are skipped.
func (s *NewsSummaryStream) Next() (NewsSummaryEvent, error) {
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") || strings.HasPrefix(line, "event:") {
			continue
		}
		payload := line
		if strings.HasPrefix(line, "data:") {
			payload = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
		if payload == "" {
			continue
		}
		var event NewsSummaryEvent
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return NewsSummaryEvent{}, errs.Wrap(errs.CodeAPI, "decoding news summary event", err)
		}
		return event, nil
	}
	if err := s.scanner.Err(); err != nil {
		return NewsSummaryEvent{}, errs.Wrap(errs.CodeTransport, "reading news summary stream", err)
	}
	return NewsSummaryEvent{}, io.EOF
}

// Close releases the underlying HTTP response body. It is safe to call more
// than once.
func (s *NewsSummaryStream) Close() error {
	if s == nil || s.body == nil {
		return nil
	}
	return s.body.Close()
}
