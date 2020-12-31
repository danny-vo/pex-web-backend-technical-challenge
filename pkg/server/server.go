package server

import (
	"github.com/danny-vo/fibonacci-backend/pkg/fibonacci"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

// Simple server wrapper onjects to contain all basic dependencies
type Server struct {
	f_sequence *fibonacci.Fibonacci
	router     *httprouter.Router
	rdb        *redis.Client
}

// Public function used to initialize an instance of Server
func Initialize_Server() *Server {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	s := &Server{
		f_sequence: fibonacci.Initialize_Fibonacci(rdb),
		router:     httprouter.New(),
		rdb:        rdb,
	}

	s.routes()
	return s
}

// Getter function for the router
func (s *Server) Get_Router() *httprouter.Router {
	return s.router
}
