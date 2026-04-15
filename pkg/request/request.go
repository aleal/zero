package request

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/aleal/zero/pkg/parser"
)

// Handler is a function that handles a request
type Handler func(w http.ResponseWriter, r *http.Request)

// GetPathParam gets a path parameter
func GetPathParam(r *http.Request, key string) string {
	return GetPathParamOrDefault(r, key, "")
}

// GetPathParamOrDefault gets a path parameter with default value
func GetPathParamOrDefault(r *http.Request, key, defaultValue string) string {
	if value := r.PathValue(key); value != "" {
		return parser.SanitizeString(value)
	}
	return defaultValue
}

// GetQueryParam gets a query parameter
func GetQueryParam(r *http.Request, key string) string {
	return GetQueryParamOrDefault(r, key, "")
}

// GetQueryParamOrDefault gets a query parameter with default value
func GetQueryParamOrDefault(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return parser.SanitizeString(value)
	}
	return defaultValue
}

func GetParsedQueryParamOrDefault[T parser.ParseType](r *http.Request, key string, defaultValue T) T {
	var value T
	err := parser.ParseString(GetQueryParam(r, key), &value)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetUploadedFile(r *http.Request, field string, maxMultipartFormSize int64) ([]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMultipartFormSize); err != nil {
		return nil, err
	}
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, fmt.Errorf("no multipart form")
	}
	files, ok := r.MultipartForm.File[field]
	if !ok || len(files) == 0 {
		return nil, fmt.Errorf("missing file in field '%s'", field)
	}
	return files, nil
}

func ReadUploadedFiles(r *http.Request, field string, maxMultipartFormSize int64) ([][]byte, error) {
	files, err := GetUploadedFile(r, field, maxMultipartFormSize)
	if err != nil {
		return nil, err
	}
	data := make([][]byte, len(files))
	for i, file := range files {
		data[i], err = readUploadedFile(file, maxMultipartFormSize)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func ReadUploadedFile(r *http.Request, field string, maxMultipartFormSize int64) ([]byte, error) {
	data, err := ReadUploadedFiles(r, field, maxMultipartFormSize)
	if err != nil {
		return nil, err
	}
	return data[0], nil
}

func readUploadedFile(file *multipart.FileHeader, maxSize int64) ([]byte, error) {
	if file.Size > maxSize {
		return nil, fmt.Errorf("file '%s' size %d exceeds limit %d", file.Filename, file.Size, maxSize)
	}
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}
