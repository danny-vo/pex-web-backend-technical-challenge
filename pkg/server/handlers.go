package server

import (
	"fmt"
	"log"
	"net/http"
)

// This function should return the current number in the Fibonacci sequence
func (s *Server) handle_current() http.HandlerFunc {
	return recovery_wrapper(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"current": %d}`, s.f_sequence.Get_Current())))
		},
	)
}

// This function should return the next number in the Fibonacci sequence and
// progress the series
func (s *Server) handle_next() http.HandlerFunc {
	return recovery_wrapper(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"next": %d}`, s.f_sequence.Get_Next(s.rdb))))
		},
	)
}

// This function returns the previous number in the Fibonacci sequence
func (s *Server) handle_previous() http.HandlerFunc {
	return recovery_wrapper(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"previous": %d}`, s.f_sequence.Get_Previous())))
		},
	)
}

// This function is simply a health check endpoint
func (s *Server) handle_health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}
}

// This function is simply a wrapper to catch occuring panics and recover gracefully
func recovery_wrapper(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		defer func() {
			if r := recover(); nil != r {
				log.Printf("Error occured: %v\n, recovered", r)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
