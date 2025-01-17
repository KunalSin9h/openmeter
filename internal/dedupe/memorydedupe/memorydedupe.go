// Package memorydedupe implements in-memory event deduplication.
package memorydedupe

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/openmeterio/openmeter/internal/dedupe"
)

const defaultSize = 1024

// Deduplicator implements in-memory event deduplication.
type Deduplicator struct {
	store *lru.Cache[string, any]
}

// NewDeduplicator returns a new {Deduplicator}.
func NewDeduplicator(size int) (*Deduplicator, error) {
	if size < 1 {
		size = defaultSize
	}

	store, err := lru.New[string, any](size)
	if err != nil {
		return nil, err
	}

	return &Deduplicator{
		store: store,
	}, nil
}

func (d *Deduplicator) IsUnique(ctx context.Context, namespace string, ev event.Event) (bool, error) {
	isContained, _ := d.store.ContainsOrAdd(dedupe.GetEventKey(namespace, ev), nil)

	return !isContained, nil
}
