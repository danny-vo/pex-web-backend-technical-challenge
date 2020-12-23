package server

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	router *httprouter.Router
}

func initialize_server() *server {
	s := &server{
		router: httprouter.New()
	}
	s.routes()
	return s
}
