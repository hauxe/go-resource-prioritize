package resourceprioritize

import (
	"context"
	"sync"
	"time"
)

// PrioritizedEntity interface
type PrioritizedEntity interface {
	GetPriority() int
	AccessResource(interface{})
}

// GreedyMutex a mutex that handles greedy entities with higher priority
type GreedyMutex struct {
	sync.Mutex
	TimeToWait time.Duration
	queue      chan PrioritizedEntity
}

// Compete compete for resource
func (mux *GreedyMutex) Compete(ctx context.Context, resource interface{}) {
	mux.Lock()
	greedyEntities := []PrioritizedEntity{}
	defer func() { go mux.Compete(ctx, resource) }()
	defer mux.Unlock()
	defer func() {
		var winner PrioritizedEntity
		for _, entity := range greedyEntities {
			if winner == nil || winner.GetPriority() < entity.GetPriority() {
				winner = entity
			}
		}
		if winner != nil {
			winner.AccessResource(resource)
		}
	}()
	// wait until have a request for resource
	entity := <-mux.queue
	greedyEntities = append(greedyEntities, entity)
	c := time.NewTicker(mux.TimeToWait)
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.C:
			c.Stop()
			return
		case entity := <-mux.queue:
			greedyEntities = append(greedyEntities, entity)
		}
	}
}
