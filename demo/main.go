package main

import (
	"fmt"
	"time"

	resourceprioritize "github.com/hauxe/go-resource-prioritize"
)

// PrioritizedEntity defines prioritized entity
type PrioritizedEntity struct {
	ID       int
	Priority int
}

// GetPriority gets identity priority
func (p *PrioritizedEntity) GetPriority() int {
	return p.Priority
}

// AccessResource access the resource
func (p *PrioritizedEntity) AccessResource(resource interface{}) {
	res, ok := resource.(map[int]int)
	if !ok {
		fmt.Println("Error resource not expected")
	}
	res[p.ID]++
}

func main() {
	sharedResource := make(map[int]int, 3)
	gateKeeper := resourceprioritize.New(100)
	gateKeeper.Start(sharedResource, time.Microsecond)
	defer gateKeeper.Stop()
	// 3 different prioritized entity
	highEntity := &PrioritizedEntity{
		ID:       1,
		Priority: 3,
	}
	mediumEntity := &PrioritizedEntity{
		ID:       2,
		Priority: 2,
	}
	lowEntity := &PrioritizedEntity{
		ID:       3,
		Priority: 1,
	}

	// start greeding resource count
	for i := 0; i < 100; i++ {
		go gateKeeper.Register(lowEntity)
		go gateKeeper.Register(mediumEntity)
		go gateKeeper.Register(highEntity)
	}
	time.Sleep(time.Second)
	fmt.Printf("Resource Status %#v\n", sharedResource)
}
