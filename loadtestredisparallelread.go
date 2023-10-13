package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	redisAddr     = "localhost:6379" // Update with your Redis server address
	redisPassword = ""               // Password if Redis requires authentication
	numRequests   = 100000           // Total number of read requests
	parallelism   = 64               // Number of parallel read goroutines
)

var (
	redisClient *redis.Client
)

func init() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
}

func readRandomKey(ctx context.Context) {
	key := fmt.Sprintf("key_%d", rand.Intn(10000)) // Adjust the key pattern to your data

	start := time.Now()
	_, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error reading key %s: %v\n", key, err)
	}
	elapsed := time.Since(start)

	fmt.Printf("Read key %s in %v\n", key, elapsed)
}

func main() {
	runtime.GOMAXPROCS(8)
	startTime := time.Now()
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numRequests/parallelism; j++ {
				readRandomKey(context.Background())
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	fmt.Printf("Total time to complete all requests: %v\n", totalTime)
}
