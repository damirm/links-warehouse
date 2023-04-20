package storage

import (
	"context"
	"net/url"
	"sync"

	"github.com/damirm/links-warehouse/internal/model"
)

type InMemoryStorage struct {
	lmu   sync.Mutex
	links map[*url.URL]*model.Link
	qmu   sync.Mutex
	queue []*url.URL
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		links: make(map[*url.URL]*model.Link),
	}
}

func (s *InMemoryStorage) SaveLink(ctx context.Context, l *model.Link) error {
	s.lmu.Lock()
	defer s.lmu.Unlock()
	s.links[l.URL] = l
	return nil
}

func (s *InMemoryStorage) EnqueueURL(ctx context.Context, u *url.URL) error {
	s.qmu.Lock()
	defer s.qmu.Unlock()
	s.queue = append(s.queue, u)
	return nil
}

func (s *InMemoryStorage) DequeueURL(context.Context) (*url.URL, error) {
	s.qmu.Lock()
	defer s.qmu.Unlock()
	if len(s.queue) == 0 {
		return nil, nil
	}
	first := s.queue[0]
	s.queue = s.queue[1:]
	return first, nil
}

func (s *InMemoryStorage) DeleteProcessedURL(ctx context.Context, du *url.URL) error {
	s.qmu.Lock()
	defer s.qmu.Unlock()
	s.queue = filter(s.queue, func(u *url.URL) bool {
		return u.String() == du.String()
	})
	return nil
}

func (s *InMemoryStorage) Transaction(ctx context.Context, fn func(context.Context, Storage) error) error {
	return fn(ctx, s)
}

func filter(lst []*url.URL, fn func(*url.URL) bool) []*url.URL {
	res := make([]*url.URL, 0, len(lst))
	for _, u := range lst {
		if fn(u) {
			res = append(res, u)
		}
	}
	return res
}
