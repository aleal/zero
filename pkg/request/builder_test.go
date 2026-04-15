package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder(context.Background(), http.MethodGet, "http://example.com")
	if b == nil {
		t.Fatal("NewBuilder returned nil")
	}
}

func TestBuilderBuildSimpleGET(t *testing.T) {
	req, err := NewBuilder(context.Background(), http.MethodGet, "http://example.com/path").
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if req.Method != http.MethodGet {
		t.Errorf("method = %q, want GET", req.Method)
	}
	if req.URL.String() != "http://example.com/path" {
		t.Errorf("url = %q, want http://example.com/path", req.URL.String())
	}
}

func TestBuilderWithMethod(t *testing.T) {
	req, err := NewBuilder(context.Background(), http.MethodGet, "http://example.com").
		WithMethod(http.MethodPost).
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if req.Method != http.MethodPost {
		t.Errorf("method = %q, want POST", req.Method)
	}
}

func TestBuilderWithURL(t *testing.T) {
	req, err := NewBuilder(context.Background(), http.MethodGet, "http://old.com").
		WithURL("http://new.com/updated").
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if req.URL.String() != "http://new.com/updated" {
		t.Errorf("url = %q, want http://new.com/updated", req.URL.String())
	}
}

func TestBuilderWithContext(t *testing.T) {
	type ctxKey string
	ctx := context.WithValue(context.Background(), ctxKey("k"), "v")
	req, err := NewBuilder(context.Background(), http.MethodGet, "http://example.com").
		WithContext(ctx).
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if req.Context().Value(ctxKey("k")) != "v" {
		t.Error("context value not propagated")
	}
}

func TestBuilderWithHeader(t *testing.T) {
	req, err := NewBuilder(context.Background(), http.MethodGet, "http://example.com").
		WithHeader("X-Custom", "hello").
		WithHeader("Authorization", "Bearer tok").
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if got := req.Header.Get("X-Custom"); got != "hello" {
		t.Errorf("X-Custom = %q, want hello", got)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer tok" {
		t.Errorf("Authorization = %q, want Bearer tok", got)
	}
}

func TestBuilderWithBodyJSON(t *testing.T) {
	payload := map[string]string{"name": "zero", "version": "1.0"}
	req, err := NewBuilder(context.Background(), http.MethodPost, "http://example.com").
		WithBodyJSON(payload).
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if ct := req.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	body, _ := io.ReadAll(req.Body)
	var decoded map[string]string
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("body unmarshal error: %v", err)
	}
	if decoded["name"] != "zero" || decoded["version"] != "1.0" {
		t.Errorf("body = %v, want {name:zero, version:1.0}", decoded)
	}
}

func TestBuilderWithBodyJSONInvalid(t *testing.T) {
	// channels can't be marshalled
	_, err := NewBuilder(context.Background(), http.MethodPost, "http://example.com").
		WithBodyJSON(make(chan int)).
		Build()
	if err == nil {
		t.Fatal("expected error for un-marshallable body, got nil")
	}
}

func TestBuilderWithFormField(t *testing.T) {
	req, err := NewBuilder(context.Background(), http.MethodPost, "http://example.com").
		WithFormField("username", "xandao").
		WithFormField("role", "admin").
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		t.Fatalf("Content-Type = %q, want multipart/form-data", ct)
	}
	// Parse the multipart form to verify fields
	_, params, err := mime.ParseMediaType(ct)
	if err != nil {
		t.Fatalf("parse media type: %v", err)
	}
	reader := multipart.NewReader(req.Body, params["boundary"])
	fields := make(map[string]string)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("reading part: %v", err)
		}
		data, _ := io.ReadAll(part)
		fields[part.FormName()] = string(data)
	}
	if fields["username"] != "xandao" {
		t.Errorf("username = %q, want xandao", fields["username"])
	}
	if fields["role"] != "admin" {
		t.Errorf("role = %q, want admin", fields["role"])
	}
}

func TestBuilderWithFormFile(t *testing.T) {
	content := []byte("file content here")
	file := NewUploadedFile(content, "test.txt", "text/plain")

	req, err := NewBuilder(context.Background(), http.MethodPost, "http://example.com").
		WithFormFile("document", file).
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		t.Fatalf("Content-Type = %q, want multipart/form-data", ct)
	}
	_, params, _ := mime.ParseMediaType(ct)
	reader := multipart.NewReader(req.Body, params["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		t.Fatalf("reading part: %v", err)
	}
	if part.FormName() != "document" {
		t.Errorf("form name = %q, want document", part.FormName())
	}
	if part.FileName() != "test.txt" {
		t.Errorf("filename = %q, want test.txt", part.FileName())
	}
	data, _ := io.ReadAll(part)
	if !bytes.Equal(data, content) {
		t.Errorf("file content = %q, want %q", data, content)
	}
}

func TestBuilderWithFormFileAndField(t *testing.T) {
	file := NewUploadedFile([]byte("img data"), "photo.png", "image/png")
	req, err := NewBuilder(context.Background(), http.MethodPost, "http://example.com").
		WithFormFile("avatar", file).
		WithFormField("description", "my photo").
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		t.Fatalf("Content-Type = %q, want multipart/form-data", ct)
	}
}

func TestBuilderFluentChaining(t *testing.T) {
	req, err := NewBuilder(context.Background(), http.MethodGet, "http://example.com").
		WithMethod(http.MethodPut).
		WithURL("http://other.com/api").
		WithHeader("Accept", "application/json").
		WithBodyJSON(map[string]int{"count": 42}).
		Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if req.Method != http.MethodPut {
		t.Errorf("method = %q, want PUT", req.Method)
	}
	if req.URL.String() != "http://other.com/api" {
		t.Errorf("url = %q", req.URL.String())
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("Accept header missing")
	}
}

func TestBuilderBuildInvalidURL(t *testing.T) {
	_, err := NewBuilder(context.Background(), http.MethodGet, "://bad-url").
		Build()
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestNewUploadedFile(t *testing.T) {
	content := []byte("hello")
	f := NewUploadedFile(content, "file.txt", "text/plain")
	if f.Name != "file.txt" {
		t.Errorf("Name = %q, want file.txt", f.Name)
	}
	if f.ContentType != "text/plain" {
		t.Errorf("ContentType = %q, want text/plain", f.ContentType)
	}
	if !bytes.Equal(f.Content, content) {
		t.Errorf("Content mismatch")
	}
}

// -- File upload helpers (GetUploadedFile, ReadUploadedFiles, ReadUploadedFile) --

func buildMultipartRequest(t *testing.T, fieldName string, files map[string][]byte) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for name, content := range files {
		part, err := w.CreateFormFile(fieldName, name)
		if err != nil {
			t.Fatalf("CreateFormFile: %v", err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestGetUploadedFile(t *testing.T) {
	req := buildMultipartRequest(t, "file", map[string][]byte{
		"a.txt": []byte("aaa"),
	})
	headers, err := GetUploadedFile(req, "file", 10<<20)
	if err != nil {
		t.Fatalf("GetUploadedFile error: %v", err)
	}
	if len(headers) != 1 {
		t.Fatalf("got %d files, want 1", len(headers))
	}
	if headers[0].Filename != "a.txt" {
		t.Errorf("filename = %q, want a.txt", headers[0].Filename)
	}
}

func TestGetUploadedFileMissingField(t *testing.T) {
	req := buildMultipartRequest(t, "file", map[string][]byte{
		"a.txt": []byte("aaa"),
	})
	_, err := GetUploadedFile(req, "other_field", 10<<20)
	if err == nil {
		t.Fatal("expected error for missing field, got nil")
	}
}

func TestGetUploadedFileNotMultipart(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader("plain body"))
	req.Header.Set("Content-Type", "text/plain")
	_, err := GetUploadedFile(req, "file", 10<<20)
	if err == nil {
		t.Fatal("expected error for non-multipart request, got nil")
	}
}

func TestReadUploadedFile(t *testing.T) {
	content := []byte("file-content-123")
	req := buildMultipartRequest(t, "doc", map[string][]byte{
		"doc.pdf": content,
	})
	data, err := ReadUploadedFile(req, "doc", 10<<20)
	if err != nil {
		t.Fatalf("ReadUploadedFile error: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("content = %q, want %q", data, content)
	}
}

func TestReadUploadedFiles(t *testing.T) {
	req := buildMultipartRequest(t, "files", map[string][]byte{
		"one.txt": []byte("one"),
		"two.txt": []byte("two"),
	})
	data, err := ReadUploadedFiles(req, "files", 10<<20)
	if err != nil {
		t.Fatalf("ReadUploadedFiles error: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("got %d files, want 2", len(data))
	}
}

func TestReadUploadedFileMissingField(t *testing.T) {
	req := buildMultipartRequest(t, "file", map[string][]byte{
		"a.txt": []byte("aaa"),
	})
	_, err := ReadUploadedFile(req, "missing", 10<<20)
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}
