package fibonacci

import (
	"sync"
	"testing"
)

func TestFibonacci_threadSafe(t *testing.T) {
	f := &Fibonacci{
		current:  0,
		next:     1,
		previous: 0,
		rwMutex:  &sync.RWMutex{},
	}

	for i := 0; i < 1000; i++ {
		go func() {
			f.GetCurrent()
			f.GetNext(mockRdb{})
			f.GetPrevious()
		}()
	}
}
