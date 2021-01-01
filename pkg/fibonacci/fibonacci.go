package fibonacci

import (
	"context"
	"errors"
	"log"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	redisFibKey = "fibonacci_current"
)

// RedisClient -
// Wrapper interface for redis client Get and Set
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// Fibonacci - Simple wrapper for the state of a Fibonacci sequence
type Fibonacci struct {
	current  *big.Int
	next     *big.Int
	previous *big.Int
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

	currentBig := new(big.Int)
	currentBig, ok := currentBig.SetString(current, 10)
	if !ok {
		return nil, errors.New("NAN retrieved from redis")
	}

	divisor := new(big.Float).SetFloat64((1 + math.Sqrt(5)) / 2.0)
	prev := new(big.Float).Quo(
		new(big.Float).SetInt(currentBig),
		divisor,
	)
	prevInt, _ := prev.Int(nil)

	return &Fibonacci{
		current:  currentBig,
		next:     new(big.Int).Add(currentBig, prevInt),
		previous: prevInt,
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
		current:  big.NewInt(0),
		next:     big.NewInt(1),
		previous: big.NewInt(0),
		rwMutex:  &sync.RWMutex{},
	}
}

// GetCurrent -
// This function will retrieve the value the sequence is currently on.
// It will also set a reading lock.
func (f *Fibonacci) GetCurrent() *big.Int {
	f.rwMutex.RLock()
	defer f.rwMutex.RUnlock()

	return f.current
}

// GetNext -
// This function will both retrieve the next value in the sequence and update
// the previous and current values.
// This function is locked from starting while any other R/W operations are occuring
func (f *Fibonacci) GetNext(rdb RedisClient) *big.Int {
	f.rwMutex.Lock()
	defer f.rwMutex.Unlock()

	oldNext := f.next
	f.previous = f.current
	f.current = f.next
	f.next = new(big.Int).Add(f.current, f.previous)

	// Store in cache to restore from in case container goes boom
	go func() {
		log.Printf(
			"Updating redis with state:\n\tprevious: %v\n\tcurrent: %v\n\tnext: %v\n\n",
			f.previous.String(), f.current.String(), f.next.String(),
		)

		if err := rdb.Set(
			context.Background(), "fibonacci_current", f.current.String(), 0,
		).Err(); nil != err {
			log.Printf("Error updating redis state: %v", err)
		}
	}()

	return oldNext
}

// GetPrevious -
// This function will retrieve the previous value in the sequence
// It will also set a reading lock
func (f *Fibonacci) GetPrevious() *big.Int {
	f.rwMutex.RLock()
	defer f.rwMutex.RUnlock()

	return f.previous
}
