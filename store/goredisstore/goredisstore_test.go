package goredisstore_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/goredisstore"
	"github.com/throttled/throttled/v2/store/storetest"
)

const (
	redisTestDB     = 1
	redisTestPrefix = "throttled-go-redis:"
)

// Demonstrates that how to initialize a RateLimiter with redis
// using go-redis library.
func ExampleNew() {
	// import "github.com/go-redis/redis/v8"

	// Initialize a redis client using go-redis
	client := redis.NewClient(&redis.Options{
		PoolSize:    10, // default
		IdleTimeout: 30 * time.Second,
		Addr:        "localhost:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
	})

	// Setup store
	store, err := goredisstore.New(client, "throttled:")
	if err != nil {
		log.Fatal(err)
	}

	// Setup quota
	quota := throttled.RateQuota{MaxRate: throttled.PerMin(20), MaxBurst: 5}

	// Then, use store and quota as arguments for NewGCRARateLimiter()
	throttled.NewGCRARateLimiter(store, quota)
}

func TestRedisStore(t *testing.T) {
	ctx := context.Background()
	c, st := setupRedis(ctx, t, 0)
	defer c.Close()
	defer clearRedis(ctx, c)

	clearRedis(ctx, c)
	storetest.TestGCRAStore(t, st)
	storetest.TestGCRAStoreTTL(t, st)
}

func BenchmarkRedisStore(b *testing.B) {
  ctx := context.Background()
	c, st := setupRedis(ctx, b, 0)
	defer c.Close()
	defer clearRedis(ctx, c)

	storetest.BenchmarkGCRAStore(ctx, b, st)
}

func clearRedis(ctx context.Context, c *redis.Client) error {
	keys, err := c.Keys(ctx, redisTestPrefix+"*").Result()
	if err != nil {
		return err
	}

	return c.Del(ctx, keys...).Err()
}

func setupRedis(ctx context.Context, tb testing.TB, ttl time.Duration) (*redis.Client, *goredisstore.GoRedisStore) {
	client := redis.NewClient(&redis.Options{
		PoolSize:    10, // default
		IdleTimeout: 30 * time.Second,
		Addr:        "localhost:6379",
		Password:    "",          // no password set
		DB:          redisTestDB, // use default DB
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		tb.Skip("redis server not available on localhost port 6379")
	}

	st, err := goredisstore.New(client, redisTestPrefix)
	if err != nil {
		client.Close()
		tb.Fatal(err)
	}

	return client, st
}
