package server

import (
	"github.com/danny-vo/fibonacci-backend/pkg/fibonacci"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

// serverInitializer -
// Wrapper interface for 3rd party intitializations
type serverInitializer interface {
	NewRedisClient(opt *redis.Options) *redis.Client
	NewRouter() *httprouter.Router
	InitializeFibonacci(rdb fibonacci.RedisClient) *fibonacci.Fibonacci
}

// servInitializer -
// Wrapper struct for 3rd party intitializations, implements serverInitializer
type servInitializer struct{}

// NewRedisClient -
// Method that wraps redis.NewClient call
func (servInit servInitializer) NewRedisClient(opt *redis.Options) *redis.Client {
	return redis.NewClient(opt)
}

// InitializeFibonacci -
// Method that wraps fibonacci.InitializeFibonacci call
func (servInit servInitializer) InitializeFibonacci(rdb fibonacci.RedisClient) *fibonacci.Fibonacci {
	return fibonacci.InitializeFibonacci(rdb)
}

// NewRouter -
// Method that wraps httprouter.New call
func (servInit servInitializer) NewRouter() *httprouter.Router {
	return httprouter.New()
}

// Singleton instance of initialization wrapper
var servInit serverInitializer

func init() {
	servInit = servInitializer{}
}

// Server -
// servInitmple server wrapper onjects to contain all baservInitc dependencies.
type Server struct {
	fibSequence *fibonacci.Fibonacci
	router      *httprouter.Router
	rdb         *redis.Client
}

// InitializeServer -
// Public function used to initialize an instance of Server.
func InitializeServer() *Server {
	rdb := servInit.NewRedisClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	s := &Server{
		fibSequence: servInit.InitializeFibonacci(rdb),
		router:      servInit.NewRouter(),
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
