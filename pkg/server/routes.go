package server

import "net/http"

// This funciton initializes all the methods, routes, and assigns handlers for
// the server's router
func (s *Server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/current", recoveryWrapper(s.handleCurrent()))
	s.router.HandlerFunc(http.MethodGet, "/next", recoveryWrapper(s.handleNext()))
	s.router.HandlerFunc(http.MethodGet, "/previous", recoveryWrapper(s.handlePrevious()))
	s.router.HandlerFunc(http.MethodGet, "/health", s.handleHealth())
}
