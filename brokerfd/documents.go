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
	"net/url"
)

const (
	pathDocuments        = "/broker-fd/documents"
	pathDocumentDetail   = "/broker-fd/documents/detail"
	pathDocumentUpload   = "/broker-fd/documents/upload"
	pathDocumentDownload = "/broker-fd/documents/download"
)

type Document struct {
	DocumentID   string `json:"document_id"`
	DocumentType string `json:"document_type"`
	FileName     string `json:"file_name"`
	FileSize     string `json:"file_size"`
	Status       string `json:"status"`
	UploadTime   string `json:"upload_time"`
}

type UploadDocumentRequest struct {
	DocumentType string `json:"document_type"`
	FileName     string `json:"file_name"`
	Content      string `json:"content"`
}

func (c *Client) UploadDocument(ctx context.Context, req UploadDocumentRequest) (*Document, error) {
	var out Document
	if err := c.post(ctx, pathDocumentUpload, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DownloadDocument(ctx context.Context, documentID string) ([]byte, error) {
	q := url.Values{}
	q.Set("document_id", documentID)
	var out []byte
	if err := c.get(ctx, pathDocumentDownload, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ListDocuments(ctx context.Context, accountID string) ([]Document, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []Document
	if err := c.get(ctx, pathDocuments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetDocumentDetail(ctx context.Context, documentID string) (*Document, error) {
	q := url.Values{}
	q.Set("document_id", documentID)
	var out Document
	if err := c.get(ctx, pathDocumentDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
