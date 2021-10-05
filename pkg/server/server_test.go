package server

import (
	"reflect"
	"testing"

	"github.com/dvo-dev/fibonacci-backend/pkg/fibonacci"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

type mockServerInitializer struct {
	rdb    *redis.Client
	router *httprouter.Router
}

func (msi mockServerInitializer) NewRedisClient(opt *redis.Options) *redis.Client {
	return msi.rdb
}

func (msi mockServerInitializer) InitializeFibonacci(rdb fibonacci.RedisClient) *fibonacci.Fibonacci {
	return &fibonacci.Fibonacci{}
}

func (msi mockServerInitializer) NewRouter() *httprouter.Router {
	return msi.router
}

func TestInitializeServer(t *testing.T) {
	mockServerInit := mockServerInitializer{
		rdb: redis.NewClient(&redis.Options{
			Addr:     "redis: 2020",
			Password: "foobar",
			DB:       42,
		}),
		router: httprouter.New(),
	}

	tests := []struct {
		name string
		want *Server
	}{
		{
			name: "happy path",
			want: &Server{
				fibSequence: &fibonacci.Fibonacci{},
				router:      mockServerInit.router,
				rdb:         mockServerInit.rdb,
			},
		},
	}
	for _, tt := range tests {

		servInit = mockServerInit
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetRouter(t *testing.T) {
	type fields struct {
		router *httprouter.Router
	}
	testRouter := httprouter.New()

	tests := []struct {
		name   string
		fields fields
		want   *httprouter.Router
	}{
		{
			name: "happy path",
			fields: fields{
				router: testRouter,
			},
			want: testRouter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				router: tt.fields.router,
			}
			if got := s.GetRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.GetRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
