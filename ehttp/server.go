package ehttp

import "net"

type RequestHandler func(ctx *RequestCtx)

type RequestCtx struct {
}

type Server struct {
	Handler RequestHandler
}

func (s *Server) ListenAndServer(addr string) error {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	if _, ok := ln.(*net.TCPListener); ok {
		return nil
	}
	return s.Serve(ln)
}

func (s *Server) Server(ln net.Listener) error {

}

func ListenAndServe(addr string, handler RequestHandler) error {
	s := &Server{
		Handler: handler,
	}
	return s.ListenAndServer(addr)
}
