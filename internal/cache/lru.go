package cache

import (
	"fmt"
	"time"

	"github.com/viccon/sturdyc"

	"game-tracker/internal/igdb"
)

type Cache struct {
	client *sturdyc.Client[[]igdb.SearchCandidate]
}

func NewCache(maxEntries int, ttl time.Duration) *Cache {
	capacity := maxEntries
	numShards := 10
	evictionPercentage := 10

	client := sturdyc.New[[]igdb.SearchCandidate](
		capacity,
		numShards,
		ttl,
		evictionPercentage,
	)

	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(query string) ([]igdb.SearchCandidate, bool) {
	result, ok := c.client.Get(query)
	return result, ok
}

func (c *Cache) Set(query string, results []igdb.SearchCandidate) {
	fmt.Printf("[Cache] SET for query: %q, storing %d results\n", query, len(results))
	c.client.Set(query, results)
}
