package fibonacci

type Fibonacci struct {
	current  uint32
	next     uint32
	previous uint32
}

func Initialize_Fibonacci() *Fibonacci {
	return &Fibonacci{
		current:  0,
		next:     1,
		previous: 0,
	}
}

func (f *Fibonacci) Get_Current() uint32 {
	return f.current
}

func (f *Fibonacci) Get_Next() uint32 {
	f.previous = f.current
	f.current = f.next
	f.next = f.current + f.previous
	return f.next
}

func (f *Fibonacci) Get_Previous() uint32 {
	return f.previous
}
