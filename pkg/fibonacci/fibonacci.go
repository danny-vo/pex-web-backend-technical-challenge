package fibonacci

import (
	"context"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient -
// Wrapper interface for redis client Get and Set
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// Fibonacci - Simple wrapper for the state of a Fibonacci sequence
type Fibonacci struct {
	current  uint64
	next     uint64
	previous uint64
	rwMutex  *sync.RWMutex
}

// This function attempts to restore a fibonacci sequence state as saved in redis
// using the "current" value
func restoreFibonacci(rdb RedisClient) (*Fibonacci, error) {
	current, err := rdb.Get(context.Background(), "fibonacci_current").Result()
	if nil != err {
		log.Printf("Error attempting to restore sequence from redis: %v", err)
		return nil, err
	}

	currentUint, err := strconv.ParseUint(current, 10, 32)
	if nil != err {
		return nil, err
	}

	prev := float64(currentUint) / ((1 + math.Sqrt(5)) / 2.0)
	roundedPrev := uint64(math.Round(prev))

	return &Fibonacci{
		current:  uint64(currentUint),
		next:     uint64(currentUint) + roundedPrev,
		previous: roundedPrev,
		rwMutex:  &sync.RWMutex{},
	}, nil
}

// InitializeFibonacci -
// This function initializes the Fibonacci wrapper to the start of the sequence.
func InitializeFibonacci(rdb RedisClient) *Fibonacci {
	if fib, err := restoreFibonacci(rdb); nil == err {
		log.Printf("Successfully restoring sequence state from redis:\n\tcurrent: %v\n\tnext: %v\n\tprevious: %v\n\n", fib.current, fib.next, fib.previous)
		return fib
	}

	log.Println("Starting with a fresh sequence")
	return &Fibonacci{
		current:  0,
		next:     1,
		previous: 0,
		rwMutex:  &sync.RWMutex{},
	}
}

// GetCurrent -
// This function will retrieve the value the sequence is currently on.
// It will also set a reading lock.
func (f *Fibonacci) GetCurrent() uint64 {
	f.rwMutex.RLock()
	defer f.rwMutex.RUnlock()

	return f.current
}

// GetNext -
// This function will both retrieve the next value in the sequence and update
// the previous and current values.
// This function is locked from starting while any other R/W operations are occuring
func (f *Fibonacci) GetNext(rdb RedisClient) uint64 {
	f.rwMutex.Lock()
	defer f.rwMutex.Unlock()

	oldNext := f.next
	f.previous = f.current
	f.current = f.next
	f.next = f.current + f.previous

	// Store in cache to restore from in case container goes boom
	go func() {
		// log.Printf(
		// 	"Updating redis with state:\n\tprevious: %v\n\tcurrent: %v\n\tnext: %v\n\n",
		// 	f.previous, f.current, f.next,
		// )

		if err := rdb.Set(
			context.Background(), "fibonacci_current", f.current, 0,
		).Err(); nil != err {
			log.Printf("Error updating redis state: %v", err)
		}
	}()

	return oldNext
}

// GetPrevious -
// This function will retrieve the previous value in the sequence
// It will also set a reading lock
func (f *Fibonacci) GetPrevious() uint64 {
	f.rwMutex.RLock()
	defer f.rwMutex.RUnlock()

	return f.previous
}
