package provider

import (
	"context"
)

type File struct {
	Name    string
	Content []byte
}

type Provider interface {
	LoadExistingFiles() error

	Watch(ctx context.Context) error

	GetFile(endpoint string) ([]byte, bool)

	ListEndpoints() []string
}
