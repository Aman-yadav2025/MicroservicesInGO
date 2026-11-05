package storage

import (
    "context"
    "io"
)


type Provider interface {
    UploadFile(ctx context.Context, file io.Reader, filename string, size int64, contentType string) (string, error)
}