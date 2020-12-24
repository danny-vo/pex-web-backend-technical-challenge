package server

import "net/http"

func (s *Server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/current/", s.handle_current())
	s.router.HandlerFunc(http.MethodGet, "/next/", s.handle_next())
	s.router.HandlerFunc(http.MethodGet, "/previous/", s.handle_previous())
}
