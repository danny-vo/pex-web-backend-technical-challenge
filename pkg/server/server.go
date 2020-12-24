package server

import (
	"github.com/danny-vo/fibonacci-backend/pkg/fibonacci"
	"github.com/julienschmidt/httprouter"
)

// Simple server wrapper onjects to contain all basic dependencies
type Server struct {
	f_sequence *fibonacci.Fibonacci
	router     *httprouter.Router
}

// Public function used to initialize an instance of Server
func Initialize_Server() *Server {
	s := &Server{
		f_sequence: fibonacci.Initialize_Fibonacci(),
		router:     httprouter.New(),
	}
	s.routes()
	return s
}

func (s *Server) Get_Router() *httprouter.Router {
	return s.router
}
