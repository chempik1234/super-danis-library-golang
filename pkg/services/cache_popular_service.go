package services

import (
	"context"
	"fmt"
	"github.com/chempik1234/super-danis-library-golang/pkg/genericports"
	"github.com/chempik1234/super-danis-library-golang/pkg/pkgports/adapters/cache/lru"
)

// CachePopularService - cache only popular objects
//
// count uses in LRU, cache somewhere (e.g. in-memory, redis, you provide repo)
//
// doesn't care about frequency of requests, only about amount of ones
//
//	type DanisService struct {
//	  cacheService CachePopularService[string, *models.Danis]
//	  storageRepository DanisStoragePort
//	}
//
//	func (s *DanisService) Get(ctx context.Context, id string) *models.Danis {
//	  obj, err := s.cacheService.Get(ctx context.Context, id)
//	  if obj != nil {
//	    return obj
//	  }
//
//	  obj2, found, err := s.storageRepository.Get(id)
//	  if found {
//	    err = s.cacheService.UpdatePopularity(ctx context.Context, obj2, 1)  // caches if object is popular
//	  }
//	  return obj
//	}
type CachePopularService[K comparable, V genericports.ObjectWithIdentifier[K]] struct {
	usesCountLRUCache *lru.CacheLRUInMemory[K, int]
	cacheStorage      genericports.GenericCachePort[K, V]

	minUses int
}

// NewCachePopularService - create new CachePopularService
//
// minUses - at what use count we cache object
//
// lruCapacity - cap for uses cache (uses LRU)
//
// cacheStorage - your repo for storing values (e.g. redis impl)
func NewCachePopularService[K comparable, V genericports.ObjectWithIdentifier[K]](minUses int, lruCapacity int, cacheStorage genericports.GenericCachePort[K, V]) *CachePopularService[K, V] {
	return &CachePopularService[K, V]{
		usesCountLRUCache: lru.NewCacheLRUInMemory[K, int](lruCapacity),
		cacheStorage:      cacheStorage,
		minUses:           minUses,
	}
}

// UpdatePopularity - updates object uses count - if more than threshold, save obj in cache
func (s *CachePopularService[K, V]) UpdatePopularity(ctx context.Context, object V, uses int) error {
	// step 1. Get current count and increase it
	// step 2. Save if popular, even if already saved - it might have been erased

	// step 1.1 - Get current uses count
	objectID := object.GetUniqueIdentifier()

	count, found, err := s.usesCountLRUCache.Get(ctx, objectID)
	if err != nil {
		return fmt.Errorf("error getting uses count: %w", err)
	}
	if !found {
		count = 0
	}

	// step 1.2 - Increase uses count
	count += uses

	err = s.usesCountLRUCache.Set(ctx, objectID, count)
	if err != nil {
		return fmt.Errorf("error updating uses count: %w", err)
	}

	// step 2. Save
	if count >= s.minUses {
		err = s.save(ctx, object)
		if err != nil {
			return fmt.Errorf("error saving in cache storage: %w", err)
		}
	}

	return nil
}

func (s *CachePopularService[K, V]) save(ctx context.Context, object V) error {
	_, err := s.cacheStorage.SaveObject(ctx, &object)
	if err != nil {
		return fmt.Errorf("error saving cache: %w", err)
	}
	return nil
}

// Get - try to get object from cache, as usual
func (s *CachePopularService[K, V]) Get(ctx context.Context, objectID K) (*V, error) {
	return s.cacheStorage.GetObjectByID(ctx, objectID)
}

// ForceSave - save regardless of uses count in LRU.
func (s *CachePopularService[K, V]) ForceSave(ctx context.Context, object V) error {
	return s.save(ctx, object)
}

// MinUsesBeforeCaching - getter for minUses
func (s *CachePopularService[K, V]) MinUsesBeforeCaching() int {
	return s.minUses
}
