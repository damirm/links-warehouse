package warehouse

import (
	"context"
	"log"
	"net/url"
)

type WarehouseService struct {
	storage Storage
}

func NewWarehouseService(storage Storage) *WarehouseService {
	return &WarehouseService{
		storage: storage,
	}
}

func (s *WarehouseService) QueueLink(link *url.URL) error {
	ctx := context.Background()
	exists, err := s.storage.LinkExists(ctx, link)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.storage.EnqueueURL(ctx, link)
}

func (s *WarehouseService) NextProcessingLink() (*url.URL, error) {
	ctx := context.Background()
	u, err := s.storage.DequeueURL(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *WarehouseService) FinishLinkProcessing(ctx context.Context, link *Link) error {
	return s.storage.Transaction(ctx, func(ctx context.Context, storage Storage) error {
		err := storage.SaveLink(ctx, link)
		if err != nil {
			log.Printf("failed to save link: %v", err)
			return err
		}

		err = storage.DeleteProcessedURL(ctx, link.URL)
		if err != nil {
			log.Printf("failed to save link: %v", err)
			return err
		}

		return nil
	})
}
