package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

type Builder interface {
	WithContext(ctx context.Context) Builder
	WithMethod(method string) Builder
	WithURL(url string) Builder
	WithHeader(key, value string) Builder
	WithFormFile(key string, file *UploadedFile) Builder
	WithFormField(key, value string) Builder
	WithBodyJSON(body any) Builder
	Build() (*http.Request, error)
}

type UploadedFile struct {
	Content     []byte
	Name        string
	ContentType string
}

func NewUploadedFile(content []byte, name string, contentType string) *UploadedFile {
	return &UploadedFile{Content: content, Name: name, ContentType: contentType}
}

type builder struct {
	ctx        context.Context
	method     string
	url        string
	headers    map[string]string
	formFiles  map[string]*UploadedFile
	formFields map[string]string
	bodyJSON   any
}

func NewBuilder(ctx context.Context, method string, url string) Builder {
	return &builder{
		ctx:        ctx,
		method:     method,
		url:        url,
		headers:    make(map[string]string),
		formFiles:  make(map[string]*UploadedFile),
		formFields: make(map[string]string),
		bodyJSON:   nil,
	}
}

func (b *builder) Build() (*http.Request, error) {
	var body bytes.Buffer
	err := b.buildBody(&body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(b.ctx, b.method, b.url, &body)
	if err != nil {
		return nil, err
	}
	b.buildHeaders(req)
	return req, nil
}

func (b *builder) WithContext(ctx context.Context) Builder {
	b.ctx = ctx
	return b
}

func (b *builder) WithMethod(method string) Builder {
	b.method = method
	return b
}

func (b *builder) WithURL(url string) Builder {
	b.url = url
	return b
}

func (b *builder) WithHeader(key, value string) Builder {
	b.headers[key] = value
	return b
}

func (b *builder) WithFormFile(key string, file *UploadedFile) Builder {
	b.formFiles[key] = file
	return b
}

func (b *builder) WithFormField(key, value string) Builder {
	b.formFields[key] = value
	return b
}

func (b *builder) WithBodyJSON(body any) Builder {
	b.bodyJSON = body
	return b
}

func (b *builder) buildForm(body *bytes.Buffer) error {
	if len(b.formFiles) > 0 || len(b.formFields) > 0 {
		writer := multipart.NewWriter(body)
		for key, file := range b.formFiles {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, file.Name))
			h.Set("Content-Type", file.ContentType)
			h.Set("Content-Length", strconv.Itoa(len(file.Content)))
			w, err := writer.CreatePart(h)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, bytes.NewReader(file.Content))
			if err != nil {
				return err
			}
		}
		for key, value := range b.formFields {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}
		err := writer.Close()
		if err != nil {
			return err
		}
		b.headers["Content-Type"] = writer.FormDataContentType()
	}
	return nil
}

func (b *builder) buildBodyJSON(body *bytes.Buffer) error {
	if b.bodyJSON != nil {
		bodyJSON, err := json.Marshal(b.bodyJSON)
		if err != nil {
			return err
		}
		body.Write(bodyJSON)
		b.headers["Content-Type"] = "application/json"
	}
	return nil
}

func (b *builder) buildBody(body *bytes.Buffer) error {
	hasForm := len(b.formFiles) > 0 || len(b.formFields) > 0
	hasJSON := b.bodyJSON != nil
	if hasForm && hasJSON {
		return fmt.Errorf("cannot use both form data and JSON body")
	}
	if hasForm {
		return b.buildForm(body)
	}
	if hasJSON {
		return b.buildBodyJSON(body)
	}
	return nil
}

func (b *builder) buildHeaders(req *http.Request) {
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}
}
