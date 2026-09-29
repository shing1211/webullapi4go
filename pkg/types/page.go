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

package types

// Page is one page of a cursor-paginated list response.
//
// Fourteen list endpoints across the data and brokerfd packages document their 200
// response as this shape rather than as a bare array:
//
//	{"data": [ ... ], "pagination_key": "eyJ2IjoxLCJsYXN0SWQiOi..."}
//
// Those methods used to decode straight into a slice, which cannot read this body at
// all: `encoding/json` reports `cannot unmarshal object into Go value of type []T` and
// the call fails. That is a loud failure rather than a silent zero, so it is not the
// class of defect the response-contract gate was built to catch silently, but it means
// the methods could not succeed against a conforming server. They return a Page now.
//
// A Page is one response, not the whole collection. Pagination is cursor-based: send
// [Page.PaginationKey] back on the next request to get the following page, and stop
// when the returned key is empty. There is no total count and no page number on the
// wire, so there is no way to know the count in advance other than by paging to the
// end.
//
// The generic parameter is the element type, so a caller reads out.Data for the rows and
// out.PaginationKey for the cursor:
//
//	page, err := c.ListFDAccounts(ctx, "")
//	for page.PaginationKey != "" {
//	    // ... consume page.Data ...
//	    page, err = c.ListFDAccounts(ctx, page.PaginationKey)
//	}
type Page[T any] struct {
	// Data is the page of results, empty when the page is past the end.
	//
	// A server that returns a present but empty data array decodes to an empty
	// non-nil slice; a server that omits the key entirely decodes to nil. Both are
	// len-zero, so test with len rather than against nil.
	Data []T `json:"data"`
	// PaginationKey is the cursor for the next page, empty when this is the last page.
	//
	// The value is opaque and server-issued: it is a base64 cursor, not an offset, so
	// it cannot be constructed, ordered, or compared by a caller. Pass it back
	// unmodified.
	PaginationKey string `json:"pagination_key"`
}
