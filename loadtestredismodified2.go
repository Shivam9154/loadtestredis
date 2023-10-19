package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	keyPrefix   = "session:"
	numUsers    = 10000000
	maxSessions = 3
	numWriters  = 100
)

func main() {
	// Initialize a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Redis database index
	})

	// Start timing
	startTime := time.Now()

	// Insert records into Redis in parallel
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 1; i <= numWriters; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()

			for j := 1; j <= numUsers/numWriters; j++ {
				userID := (writerID-1)*(numUsers/numWriters) + j
				numSessions := maxSessions

				// Calculate the number of sessions based on user ID
				if userID > 5000000 {
					numSessions = 2
				}
				if userID > 8000000 {
					numSessions = 1
				}

				for k := 0; k < numSessions; k++ {
					sessionID := keyPrefix + strconv.Itoa(userID) + ":session" + strconv.Itoa(k)
					expiryTime := time.Now().Add(time.Duration(rand.Intn(3600)) * time.Second) // Random expiry time within 1 hour

					// Use HMSET to insert session data into Redis hash
					_, err := rdb.HMSet(ctx, sessionID, map[string]interface{}{
						"sessionid":   sessionID,
						"userid":      "user:" + strconv.Itoa(userID),
						"token":       generateRandomToken(),
						"expiry_time": expiryTime.Format(time.RFC3339),
					}).Result()

					if err != nil {
						fmt.Printf("Error inserting record %s: %v\n", sessionID, err)
						return
					}
				}
			}
		}(i)
	}

	// Wait for all parallel writers to finish
	wg.Wait()

	// End timing
	endTime := time.Now()

	// Calculate total time taken
	totalTime := endTime.Sub(startTime)

	// Print results
	fmt.Printf("Total Time to Insert: %.2f seconds\n", totalTime.Seconds())
}

func generateRandomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
