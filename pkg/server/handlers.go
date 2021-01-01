package server

import (
	"fmt"
	"log"
	"net/http"
)

// handleCurrent -
// This function should return the current number in the Fibonacci sequence.
func (s *Server) handleCurrent() http.HandlerFunc {
	return recoveryWrapper(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"current": %d}`, s.fibSequence.GetCurrent())))
		},
	)
}

// handleNext -
// This function should return the next number in the Fibonacci sequence and
// progress the series.
func (s *Server) handleNext() http.HandlerFunc {
	return recoveryWrapper(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"next": %d}`, s.fibSequence.GetNext(s.rdb))))
		},
	)
}

// handlePrevious -
// This function returns the previous number in the Fibonacci sequence.
func (s *Server) handlePrevious() http.HandlerFunc {
	return recoveryWrapper(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"previous": %d}`, s.fibSequence.GetPrevious())))
		},
	)
}

// handleHealth -
// This function is simply a health check endpoint.
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}
}

// This function is simply a wrapper to catch occuring panics and recover gracefully
func recoveryWrapper(h http.HandlerFunc) http.HandlerFunc {
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
