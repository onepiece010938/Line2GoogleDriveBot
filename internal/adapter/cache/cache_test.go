package cache

import (
	"context"
	"os"
	"testing"
)

var testCache *Cache

func TestMain(m *testing.M) {
	rootCtx, rootCtxCancelFunc := context.WithCancel(context.Background())

	testCache = NewCache(InitBigCache(rootCtx))
	defer testCache.cache.Close()
	defer rootCtxCancelFunc()
	os.Exit(m.Run())
}
