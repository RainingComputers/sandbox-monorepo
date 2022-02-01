package storage

import (
	"context"
	"examples/bloggy/pkg/models"
)

type MemoryStore struct {
	mp map[string]models.Post
}

func CreateMemoryStore() Storage {
	return &MemoryStore{mp: map[string]models.Post{}}
}

func (m *MemoryStore) Insert(ctx context.Context, post models.Post) error {
	_, ok := m.mp[post.Title]

	if ok {
		return ErrAlreadyExists
	}

	m.mp[post.Title] = post

	return nil
}

func (m *MemoryStore) Find(ctx context.Context, title string) (models.Post, error) {
	foundPost, ok := m.mp[title]

	if !ok {
		return foundPost, ErrDoesNotExist
	}

	return foundPost, nil
}

func (m *MemoryStore) Remove(ctx context.Context, title string) error {
	_, ok := m.mp[title]

	if !ok {
		return ErrDoesNotExist
	}

	delete(m.mp, title)

	return nil
}

func (m *MemoryStore) Modify(ctx context.Context, title string, post models.Post) error {
	_, ok := m.mp[title]

	if !ok {
		return ErrDoesNotExist
	}

	m.mp[title] = post

	return nil
}

func (m *MemoryStore) All(_ context.Context) ([]models.Post, error) {
	var allPosts []models.Post

	for _, v := range m.mp {
		allPosts = append(allPosts, v)
	}

	return allPosts, nil
}

func (m *MemoryStore) Disconnect(_ context.Context) error {
	return nil
}

func (m *MemoryStore) Clean(_ context.Context) error {
	for k := range m.mp {
		delete(m.mp, k)
	}

	return nil
}
