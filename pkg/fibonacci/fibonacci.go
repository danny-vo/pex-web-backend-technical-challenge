package fibonacci

import "sync"

type Fibonacci struct {
	current  uint32
	next     uint32
	previous uint32
	rw_mutex *sync.RWMutex
}

// This function initializes the Fibonacci construct to the start of the sequence
func Initialize_Fibonacci() *Fibonacci {
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
func (f *Fibonacci) Get_Next() uint32 {
	f.rw_mutex.Lock()
	defer f.rw_mutex.Unlock()

	f.previous = f.current
	f.current = f.next
	f.next = f.current + f.previous

	return f.next
}

// This function will retrieve the previous value in the sequence
// It will also set a reading lock
func (f *Fibonacci) Get_Previous() uint32 {
	f.rw_mutex.RLock()
	defer f.rw_mutex.RUnlock()

	return f.previous
}
