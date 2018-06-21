package resourceprioritise

import (
	"sync"
	"time"
)

// PrioritisedEntity interface
type PrioritisedEntity interface {
	GetPriority() int
	AccessResource(interface{})
}

// GreedyMutex a mutex that handles greedy entities with higher priority
type GreedyMutex struct {
	sync.Mutex
	TimeToWait time.Duration
}

var (
	mux   *GreedyMutex
	queue chan PrioritisedEntity
)

// Compete compete for resource
func Compete(resource interface{}) {
	mux.Lock()
	greedyEntities := []PrioritisedEntity{}
	defer func() { go Compete(resource) }()
	defer mux.Unlock()
	defer func() {
		var winner PrioritisedEntity
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
	entity := <-queue
	greedyEntities = append(greedyEntities, entity)
	c := time.NewTicker(mux.TimeToWait)
	for {
		select {
		case <-c.C:
			c.Stop()
			return
		case entity := <-queue:
			greedyEntities = append(greedyEntities, entity)
		}
	}
}
