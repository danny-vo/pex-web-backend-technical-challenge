package fibonacci

import (
	"context"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Fibonacci struct {
	current  uint32
	next     uint32
	previous uint32
	rw_mutex *sync.RWMutex
}

// This function attempts to restore a fibonacci sequence state as saved in redis
// using the "current" value
func restore_fibonacci(rdb *redis.Client) (*Fibonacci, error) {
	current, err := rdb.Get(context.Background(), "fibonacci_current").Result()
	if nil != err {
		log.Printf("Error attempting to restore sequence from redis: %v", err)
		return nil, err
	}

	current_uint, err := strconv.ParseUint(current, 10, 32)
	if nil != err {
		return nil, err
	}

	prev := float64(current_uint) / ((1 + math.Sqrt(5)) / 2.0)
	rounded_prev := uint32(math.Round(prev))

	return &Fibonacci{
		current:  uint32(current_uint),
		next:     uint32(current_uint) + rounded_prev,
		previous: rounded_prev,
		rw_mutex: &sync.RWMutex{},
	}, nil
}

// This function initializes the Fibonacci construct to the start of the sequence
func Initialize_Fibonacci(rdb *redis.Client) *Fibonacci {
	if fib, err := restore_fibonacci(rdb); nil == err {
		log.Printf("Successfully restoring sequence state from redis:\n\tcurrent: %v\n\tnext: %v\n\tprevious: %v\n\n", fib.current, fib.next, fib.previous)
		return fib
	}

	log.Println("Starting with a fresh sequence")
	return &Fibonacci{
		current:  0,
		next:     1,
		previous: 0,
		rw_mutex: &sync.RWMutex{},
	}
}

// This function will retrieve the value the sequence is currently on
// It will also set a reading lock
func (f *Fibonacci) Get_Current() uint32 {
	f.rw_mutex.RLock()
	defer f.rw_mutex.RUnlock()

	return f.current
}

// This function will both retrieve the next value in the sequence and update
// the previous and current values
// This function is locked from starting while any other R/W operations are occuring
func (f *Fibonacci) Get_Next(rdb *redis.Client) uint32 {
	f.rw_mutex.Lock()
	defer f.rw_mutex.Unlock()

	f.previous = f.current
	f.current = f.next
	f.next = f.current + f.previous

	// Store in cache to restore from in case container goes boom
	go func() {
		log.Printf(
			"Updating redis with state:\n\tprevious: %v\n\tcurrent: %v\n\tnext: %v\n\n",
			f.previous, f.current, f.next,
		)

		if err := rdb.Set(
			context.Background(), "fibonacci_current", f.current, 0,
		).Err(); nil != err {
			log.Printf("Error updating redis state: %v", err)
		}
	}()

	return f.next
}

// This function will retrieve the previous value in the sequence
// It will also set a reading lock
func (f *Fibonacci) Get_Previous() uint32 {
	f.rw_mutex.RLock()
	defer f.rw_mutex.RUnlock()

	return f.previous
}
