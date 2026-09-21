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

package brokerfd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
)

const (
	testAppKey    = "test-app-key"
	testAppSecret = "test-app-secret"
)

func TestUploadDocument(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Document{
			DocumentID:   "D1",
			DocumentType: "TAX",
			FileName:     "tax_2025.pdf",
			FileSize:     "1024",
			Status:       "uploaded",
			UploadTime:   "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.UploadDocument(context.Background(), UploadDocumentRequest{
		DocumentType: "TAX",
		FileName:     "tax_2025.pdf",
		Content:      "base64content",
	})
	if err != nil {
		t.Fatalf("UploadDocument error = %v", err)
	}
	if got.DocumentID != "D1" {
		t.Fatalf("DocumentID = %s, want D1", got.DocumentID)
	}
	if capturedReq.URL.Path != pathDocumentUpload {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathDocumentUpload)
	}
	if capturedReq.Method != http.MethodPost {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodPost)
	}
}

func TestDownloadDocument(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode("ZmlsZSBjb250ZW50") // base64 of "file content"
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.DownloadDocument(context.Background(), "D1")
	if err != nil {
		t.Fatalf("DownloadDocument error = %v", err)
	}
	if string(got) != "file content" {
		t.Fatalf("content = %s, want file content", string(got))
	}
	if capturedReq.URL.Path != pathDocumentDownload {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathDocumentDownload)
	}
	if capturedReq.URL.Query().Get("document_id") != "D1" {
		t.Fatalf("document_id = %s, want D1", capturedReq.URL.Query().Get("document_id"))
	}
}

func TestListDocuments(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Document{
			{DocumentID: "D1", DocumentType: "TAX", FileName: "tax_2025.pdf", FileSize: "1024", Status: "uploaded", UploadTime: "2026-01-01"},
			{DocumentID: "D2", DocumentType: "CONTRACT", FileName: "contract.pdf", FileSize: "2048", Status: "uploaded", UploadTime: "2026-01-02"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListDocuments(context.Background(), "A1")
	if err != nil {
		t.Fatalf("ListDocuments error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathDocuments {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathDocuments)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetDocumentDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Document{
			DocumentID:   "D1",
			DocumentType: "TAX",
			FileName:     "tax_2025.pdf",
			FileSize:     "1024",
			Status:       "uploaded",
			UploadTime:   "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetDocumentDetail(context.Background(), "D1")
	if err != nil {
		t.Fatalf("GetDocumentDetail error = %v", err)
	}
	if got.DocumentID != "D1" {
		t.Fatalf("DocumentID = %s, want D1", got.DocumentID)
	}
	if capturedReq.URL.Path != pathDocumentDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathDocumentDetail)
	}
	if capturedReq.URL.Query().Get("document_id") != "D1" {
		t.Fatalf("document_id = %s, want D1", capturedReq.URL.Query().Get("document_id"))
	}
}

func TestDownloadDocumentHandlesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	_, err = c.DownloadDocument(context.Background(), "D999")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}
