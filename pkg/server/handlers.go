package server

import (
	"fmt"
	"net/http"
)

// This function should return the current number in the Fibonacci sequence
func (s *Server) handle_current() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"current": %d}`, s.f_sequence.Get_Current())))
	}
}

// This function should return the next number in the Fibonacci sequence and
// progress the series
func (s *Server) handle_next() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"next": %d}`, s.f_sequence.Get_Next())))
	}
}

// This function returns the previous number in the Fibonacci sequence
func (s *Server) handle_previous() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"previous": %d}`, s.f_sequence.Get_Previous())))
	}
}

// This function is simply a health check endpoint
func (s *Server) handle_health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-TYpe", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"`))
	}
}
