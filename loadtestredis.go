package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Initialize a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Redis database index
	})

	// Number of records to insert
	numRecords := 1000000

	// Start timing
	startTime := time.Now()

	// Insert records into Redis
	ctx := context.Background()
	for i := 0; i < numRecords; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		_, err := rdb.Set(ctx, key, value, 0).Result()
		if err != nil {
			fmt.Printf("Error inserting record %d: %v\n", i, err)
			return
		}
	}

	// End timing
	endTime := time.Now()

	// Calculate total time taken
	totalTime := endTime.Sub(startTime)

	// Calculate mean average to insert
	meanAverage := totalTime.Seconds() / float64(numRecords)

	// Print results
	fmt.Printf("Total Time to Insert: %.2f seconds\n", totalTime.Seconds())
	fmt.Printf("Mean Average to Insert: %.6f seconds per record\n", meanAverage)
}
