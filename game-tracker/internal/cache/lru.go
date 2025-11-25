package cache

import (
	"fmt"
	"time"

	"github.com/viccon/sturdyc"

	"game-tracker/internal/igdb"
)

// Cache wraps sturdyc.Client for IGDB search results
type Cache struct {
	client *sturdyc.Client[[]igdb.SearchCandidate]
}

// NewCache creates a new cache using sturdyc with specified max entries and TTL
func NewCache(maxEntries int, ttl time.Duration) *Cache {
	// Create sturdyc client with capacity and TTL
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

// Get retrieves cached results for a query if they exist and haven't expired
func (c *Cache) Get(query string) ([]igdb.SearchCandidate, bool) {
	result, ok := c.client.Get(query)
	return result, ok
}

// Set stores search results in the cache with TTL
func (c *Cache) Set(query string, results []igdb.SearchCandidate) {
	fmt.Printf("[Cache] SET for query: %q, storing %d results\n", query, len(results))
	c.client.Set(query, results)
}
