// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package redis

import (
	"context"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/zeebo/errs"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/storage"
)

var (
	// Error is a redis error
	Error = errs.Class("redis error")

	mon = monkit.Package()
)

// TODO(coyle): this should be set to 61 * time.Minute after we implement Ping and Refresh on Overlay Cache
// This disables the TTL since the Set command only includes a TTL if it is greater than 0
const defaultNodeExpiration = 0 * time.Minute

// Client is the entrypoint into Redis
type Client struct {
	db  *redis.Client
	TTL time.Duration
}

// NewClient returns a configured Client instance, verifying a successful connection to redis
func NewClient(address, password string, db int) (*Client, error) {
	client := &Client{
		db: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
		TTL: defaultNodeExpiration,
	}

	// ping here to verify we are able to connect to redis with the initialized client.
	if err := client.db.Ping().Err(); err != nil {
		return nil, Error.New("ping failed: %v", err)
	}

	return client, nil
}

// NewClientFrom returns a configured Client instance from a redis address, verifying a successful connection to redis
func NewClientFrom(address string) (*Client, error) {
	redisurl, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if redisurl.Scheme != "redis" {
		return nil, Error.New("not a redis:// formatted address")
	}

	q := redisurl.Query()

	db, err := strconv.Atoi(q.Get("db"))
	if err != nil {
		return nil, err
	}

	return NewClient(redisurl.Host, q.Get("password"), db)
}

// Get looks up the provided key from redis returning either an error or the result.
func (client *Client) Get(ctx context.Context, key storage.Key) (_ storage.Value, err error) {
	defer mon.Task()(&ctx)(&err)
	if key.IsZero() {
		return nil, storage.ErrEmptyKey.New("")
	}

	value, err := client.db.Get(string(key)).Bytes()
	if err == redis.Nil {
		return nil, storage.ErrKeyNotFound.New(key.String())
	}
	if err != nil {
		return nil, Error.New("get error: %v", err)
	}
	return value, nil
}

// Put adds a value to the provided key in redis, returning an error on failure.
func (client *Client) Put(ctx context.Context, key storage.Key, value storage.Value) (err error) {
	defer mon.Task()(&ctx)(&err)
	if key.IsZero() {
		return storage.ErrEmptyKey.New("")
	}

	err = client.db.Set(key.String(), []byte(value), client.TTL).Err()
	if err != nil {
		return Error.New("put error: %v", err)
	}
	return nil
}

// List returns either a list of keys for which boltdb has values or an error.
func (client *Client) List(ctx context.Context, first storage.Key, limit int) (_ storage.Keys, err error) {
	defer mon.Task()(&ctx)(&err)
	return storage.ListKeys(ctx, client, first, limit)
}

// Delete deletes a key/value pair from redis, for a given the key
func (client *Client) Delete(ctx context.Context, key storage.Key) (err error) {
	defer mon.Task()(&ctx)(&err)
	if key.IsZero() {
		return storage.ErrEmptyKey.New("")
	}

	err = client.db.Del(key.String()).Err()
	if err != nil {
		return Error.New("delete error: %v", err)
	}
	return nil
}

// Close closes a redis client
func (client *Client) Close() error {
	return client.db.Close()
}

// GetAll is the bulk method for gets from the redis data store.
// The maximum keys returned will be storage.LookupLimit. If more than that
// is requested, an error will be returned
func (client *Client) GetAll(ctx context.Context, keys storage.Keys) (_ storage.Values, err error) {
	defer mon.Task()(&ctx)(&err)
	if len(keys) > storage.LookupLimit {
		return nil, storage.ErrLimitExceeded
	}

	keyStrings := make([]string, len(keys))
	for i, v := range keys {
		keyStrings[i] = v.String()
	}

	results, err := client.db.MGet(keyStrings...).Result()
	if err != nil {
		return nil, err
	}
	values := []storage.Value{}
	for _, result := range results {
		if result == nil {
			values = append(values, nil)
		} else {
			s, ok := result.(string)
			if !ok {
				return nil, Error.New("invalid result type %T", result)
			}
			values = append(values, storage.Value(s))
		}

	}
	return values, nil
}

// Iterate iterates over items based on opts
func (client *Client) Iterate(ctx context.Context, opts storage.IterateOptions, fn func(context.Context, storage.Iterator) error) (err error) {
	defer mon.Task()(&ctx)(&err)
	var all storage.Items
	if !opts.Reverse {
		all, err = client.allPrefixedItems(opts.Prefix, opts.First, nil)
	} else {
		all, err = client.allPrefixedItems(opts.Prefix, nil, opts.First)
	}
	if err != nil {
		return err
	}
	if !opts.Recurse {
		all = storage.SortAndCollapse(all, opts.Prefix)
	}
	if opts.Reverse {
		all = storage.ReverseItems(all)
	}
	return fn(ctx, &storage.StaticIterator{
		Items: all,
	})
}

// FlushDB deletes all keys in the currently selected DB.
func (client *Client) FlushDB() error {
	_, err := client.db.FlushDB().Result()
	return err
}

func (client *Client) allPrefixedItems(prefix, first, last storage.Key) (storage.Items, error) {
	var all storage.Items
	seen := map[string]struct{}{}

	match := string(escapeMatch([]byte(prefix))) + "*"
	it := client.db.Scan(0, match, 0).Iterator()
	for it.Next() {
		key := it.Val()
		if !first.IsZero() && storage.Key(key).Less(first) {
			continue
		}
		if !last.IsZero() && last.Less(storage.Key(key)) {
			continue
		}

		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		value, err := client.db.Get(key).Bytes()
		if err != nil {
			return nil, err
		}

		all = append(all, storage.ListItem{
			Key:      storage.Key(key),
			Value:    storage.Value(value),
			IsPrefix: false,
		})
	}

	sort.Sort(all)

	return all, nil
}
