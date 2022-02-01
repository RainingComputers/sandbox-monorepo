package storage

import (
	"context"
	"errors"
	"examples/bloggy/pkg/models"
)

var ErrDoesNotExist = errors.New("storage: post does not exists")
var ErrAlreadyExists = errors.New("storage: post already exists")

type Storage interface {
	Insert(ctx context.Context, post models.Post) error
	Find(ctx context.Context, title string) (models.Post, error)
	Remove(ctx context.Context, title string) error
	Modify(ctx context.Context, title string, post models.Post) error
	All(ctx context.Context) ([]models.Post, error)
	Disconnect(ctx context.Context) error
	Clean(ctx context.Context) error
}
