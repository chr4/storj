// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package testsuite

import (
	"strconv"
	"sync"
	"testing"

	"storj.io/storj/storage"
)

// RunTests runs common storage.KeyValueStore tests
func RunTests(t *testing.T, store storage.KeyValueStore) {
	// store = storelogger.NewTest(t, store)

	t.Run("CRUD", func(t *testing.T) { testCRUD(t, store) })
	t.Run("Constraints", func(t *testing.T) { testConstraints(t, store) })
	t.Run("Iterate", func(t *testing.T) { testIterate(t, store) })
	t.Run("IterateAll", func(t *testing.T) { testIterateAll(t, store) })
	t.Run("Prefix", func(t *testing.T) { testPrefix(t, store) })

	t.Run("List", func(t *testing.T) { testList(t, store) })
	t.Run("ListV2", func(t *testing.T) { testListV2(t, store) })

	t.Run("Parallel", func(t *testing.T) { testParallel(t, store) })
}

func testConstraints(t *testing.T, store storage.KeyValueStore) {
	var items storage.Items
	for i := 0; i < storage.LookupLimit+5; i++ {
		items = append(items, storage.ListItem{
			Key:   storage.Key("test-" + strconv.Itoa(i)),
			Value: storage.Value("xyz"),
		})
	}

	var wg sync.WaitGroup
	for _, item := range items {
		key := item.Key
		value := item.Value
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := store.Put(ctx, key, value)
			if err != nil {
				t.Fatal("store.Put err:", err)
			}
		}()
	}
	wg.Wait()
	defer cleanupItems(store, items)

	t.Run("Put Empty", func(t *testing.T) {
		var key storage.Key
		var val storage.Value
		defer func() { _ = store.Delete(ctx, key) }()

		err := store.Put(ctx, key, val)
		if err == nil {
			t.Fatal("putting empty key should fail")
		}
	})

	t.Run("GetAll limit", func(t *testing.T) {
		_, err := store.GetAll(ctx, items[:storage.LookupLimit].GetKeys())
		if err != nil {
			t.Fatalf("GetAll LookupLimit should succeed: %v", err)
		}

		_, err = store.GetAll(ctx, items[:storage.LookupLimit+1].GetKeys())
		if err == nil && err == storage.ErrLimitExceeded {
			t.Fatalf("GetAll LookupLimit+1 should fail: %v", err)
		}
	})

	t.Run("List limit", func(t *testing.T) {
		keys, err := store.List(ctx, nil, storage.LookupLimit)
		if err != nil || len(keys) != storage.LookupLimit {
			t.Fatalf("List LookupLimit should succeed: %v / got %d", err, len(keys))
		}
		_, err = store.List(ctx, nil, storage.LookupLimit+1)
		if err != nil || len(keys) != storage.LookupLimit {
			t.Fatalf("List LookupLimit+1 shouldn't fail: %v / got %d", err, len(keys))
		}
	})
}
