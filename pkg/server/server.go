package server

import (
	"github.com/danny-vo/fibonacci-backend/pkg/fibonacci"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

// Server -
// Simple server wrapper onjects to contain all basic dependencies.
type Server struct {
	fibSequence *fibonacci.Fibonacci
	router      *httprouter.Router
	rdb         *redis.Client
}

// InitializeServer -
// Public function used to initialize an instance of Server.
func InitializeServer() *Server {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	s := &Server{
		fibSequence: fibonacci.InitializeFibonacci(rdb),
		router:      httprouter.New(),
		rdb:         rdb,
	}

	s.routes()
	return s
}

// GetRouter -
// Getter function for the router.
func (s *Server) GetRouter() *httprouter.Router {
	return s.router
}
