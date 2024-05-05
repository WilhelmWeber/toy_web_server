package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/WilhelmWeber/toy_web_server/src/final"
	"github.com/WilhelmWeber/toy_web_server/src/handler"
	"github.com/WilhelmWeber/toy_web_server/src/toy"
)

type ToyServer struct {
	host    string
	port    int
	handler *handler.Handler
}

func NewServer(host string, port int) *ToyServer {
	return &ToyServer{
		host:    host,
		port:    port,
		handler: handler.NewHadndler(),
	}
}

func (s *ToyServer) Run() {
	addr := s.host + ":" + strconv.Itoa(s.port)
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Cannot Listen Because of Err: " + err.Error())
	}
	fmt.Println("Server Listen at port: " + strconv.Itoa(s.port))

	for {
		//IGNORE ERR
		conn, _ := listner.Accept()
		go s.handler.HandleConnection(conn)
	}
}

// Wrapping
func (s *ToyServer) AddRoute(
	method final.Method,
	route string,
	controller func(c *toy.ToyContext),
) {
	s.handler.AddRoute(method, route, controller)
}

func (s *ToyServer) Get(
	route string,
	controller func(c *toy.ToyContext),
) {
	s.AddRoute(final.GET, route, controller)
}

func (s *ToyServer) Post(
	route string,
	controller func(c *toy.ToyContext),
) {
	s.AddRoute(final.POST, route, controller)
}

func (s *ToyServer) Put(
	route string,
	controller func(c *toy.ToyContext),
) {
	s.AddRoute(final.PUT, route, controller)
}

func (s *ToyServer) Delete(
	route string,
	controller func(c *toy.ToyContext),
) {
	s.AddRoute(final.DELETE, route, controller)
}
