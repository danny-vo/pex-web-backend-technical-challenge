package server

import "net/http"

// This funciton initializes all the methods, routes, and assigns handlers for
// the server's router
func (s *Server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/current/", s.handle_current())
	s.router.HandlerFunc(http.MethodGet, "/next/", s.handle_next())
	s.router.HandlerFunc(http.MethodGet, "/previous/", s.handle_previous())
	s.router.HandlerFunc(http.MethodGet, "/health/", s.handle_health())
}
