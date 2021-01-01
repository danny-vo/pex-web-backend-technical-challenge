package server

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
)

// fibonacciSequence -
// Simple wrapper interface for accessing Server's fibSequence to make
// testing easier.
type fibonacciSequence interface {
	GetCurrent(s *Server) *big.Int
	GetNext(s *Server) *big.Int
	GetPrevious(s *Server) *big.Int
}

// fibonacciSeq -
// Implements fibonacciSequence and wraps Server fibSequence access
type fibonacciSeq struct{}

// GetCurrent -
// This method retrieves the given Server's current fibonacci number
func (fs fibonacciSeq) GetCurrent(s *Server) *big.Int {
	return s.fibSequence.GetCurrent()
}

// GetNext -
// This method retrieves the given Server's next fibonacci number
func (fs fibonacciSeq) GetNext(s *Server) *big.Int {
	return s.fibSequence.GetNext(s.rdb)
}

// GetPrevious -
// This method retrieves the given Server's previous fibonacci number
func (fs fibonacciSeq) GetPrevious(s *Server) *big.Int {
	return s.fibSequence.GetPrevious()
}

var fibSeq fibonacciSequence

func init() {
	fibSeq = fibonacciSeq{}
}

// handleCurrent -
// This function should return the current number in the Fibonacci sequence.
func (s *Server) handleCurrent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"current": %v}`, fibSeq.GetCurrent(s))))
	}
}

// handleNext -
// This function should return the next number in the Fibonacci sequence and
// progress the series.
func (s *Server) handleNext() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"next": %v}`, fibSeq.GetNext(s))))
	}
}

// handlePrevious -
// This function returns the previous number in the Fibonacci sequence.
func (s *Server) handlePrevious() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"previous": %v}`, fibSeq.GetPrevious(s))))
	}
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
