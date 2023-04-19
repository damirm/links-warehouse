package storage

import "context"

type NoopMigrator struct {
}

func (m *NoopMigrator) Migrate(context.Context) error {
	return nil
}

type InMemoryStorage struct {
}
