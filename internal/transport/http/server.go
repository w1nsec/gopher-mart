package httpserver

import (
	"net"
	"net/http"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(address string, router http.Handler) (srv *HTTPServer, err error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{Server: &http.Server{
		Addr:    addr.String(),
		Handler: router,
	}}, nil

}

func (srv *HTTPServer) Run() error {
	return srv.ListenAndServe()
}
