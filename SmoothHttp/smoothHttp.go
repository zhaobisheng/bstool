package smoothHttp

import (
	"net/http"
	"time"
)

const (
	DEFAULT_READ_TIMEOUT  = 60 * time.Second
	DEFAULT_WRITE_TIMEOUT = DEFAULT_READ_TIMEOUT
)

// refer http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT, false).ListenAndServe()
}

func ListenAndServeReturn(addr string, handler http.Handler, restart bool) (*Server, error) {
	myServer := NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT, restart)
	err := myServer.ListenAndServe()
	return myServer, err
}

// refer http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT, false).ListenAndServeTLS(certFile, keyFile)
}
