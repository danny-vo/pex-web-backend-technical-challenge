package fibonacci

import "sync"

type Fibonacci struct {
	current  uint32
	next     uint32
	previous uint32
	rw_mutex *sync.RWMutex
}

func Initialize_Fibonacci() *Fibonacci {
	return &Fibonacci{
		current:  0,
		next:     1,
		previous: 0,
		rw_mutex: &sync.RWMutex{},
	}
}

func (f *Fibonacci) Get_Current() uint32 {
	f.rw_mutex.RLock()
	defer f.rw_mutex.RUnlock()

	return f.current
}

func (f *Fibonacci) Get_Next() uint32 {
	f.rw_mutex.Lock()
	defer f.rw_mutex.Unlock()

	f.previous = f.current
	f.current = f.next
	f.next = f.current + f.previous

	return f.next
}

func (f *Fibonacci) Get_Previous() uint32 {
	f.rw_mutex.RLock()
	defer f.rw_mutex.RUnlock()

	return f.previous
}
