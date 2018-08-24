package lru_test

import (
	"testing"

	"github.com/cjsaylor/goutil/lru"
)

func TestSet(t *testing.T) {
	evictionEntries := make(map[string]string)
	cache := lru.NewCache(1, func(key, value interface{}) {
		evictionEntries[key.(string)] = value.(string)
	})
	cache.Set("a", "foo")
	cache.Set("b", "bar")
	if _, ok := evictionEntries["a"]; !ok {
		t.Error("expected eviction of 'a', but didn't find it")
	}
}

func TestSetOverwriteAndBumpRecentlyUsed(t *testing.T) {
	cache := lru.NewCache(2, func(key, value interface{}) {})
	cache.Set("a", "foo")
	cache.Set("b", "b")
	cache.Set("a", "bar")
	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to have been evicted after resetting")
	}
}

func TestGet(t *testing.T) {
	evictions := []string{}
	cache := lru.NewCache(2, func(key, value interface{}) {
		evictions = append(evictions, key.(string))
	})
	cache.Set("a", "foo")
	cache.Set("b", "foo")
	if _, ok := cache.Get("a"); !ok {
		t.Error("Expected 'a' to still exist in the cache")
	}
	cache.Set("c", "foo")
	if evictions[0] != "b" {
		t.Error("Expected 'b' to be evicted because 'a' was recently used")
	}
	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to be removed")
	}
	for _, key := range []string{"a", "c"} {
		if _, ok := cache.Get(key); !ok {
			t.Errorf("Expected %v to still be in the cache", key)
		}
	}
}

func TestRemove(t *testing.T) {
	cache := lru.NewCache(1, func(key, value interface{}) {})
	cache.Set("a", "foo")
	if res, ok := cache.Remove("a"); res.(string) != "foo" || !ok {
		t.Error("Expected to return value removed")
	}
	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be removed")
	}
}

func TestRemoveOldest(t *testing.T) {
	cache := lru.NewCache(3, func(key, value interface{}) {})
	cache.Set("a", "a")
	cache.Set("b", "b")
	cache.Set("c", "c")
	if res, ok := cache.RemoveOldest(); res.(string) != "a" || !ok {
		t.Error("Expected to return 'a' entry after removal")
	}
}
