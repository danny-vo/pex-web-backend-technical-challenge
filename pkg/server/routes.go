package server

import "net/http"

// This funciton initializes all the methods, routes, and assigns handlers for
// the server's router
func (s *Server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/current", s.handleCurrent())
	s.router.HandlerFunc(http.MethodGet, "/next", s.handleNext())
	s.router.HandlerFunc(http.MethodGet, "/previous", s.handlePrevious())
	s.router.HandlerFunc(http.MethodGet, "/health", s.handleHealth())
}
