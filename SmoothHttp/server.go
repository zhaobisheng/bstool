package smoothHttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD    = 3
)

// HTTP server that supported graceful shutdown or restart
type Server struct {
	HttpServer   *http.Server
	Listener     net.Listener
	isGraceful   bool
	signalChan   chan os.Signal
	shutdownChan chan bool
	restartFlag  bool
	LinsterList  []net.Listener
	isPPROF      bool
}

func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration, canRestart bool) *Server {
	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true
	}

	return &Server{
		HttpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		restartFlag:  canRestart,
		isGraceful:   isGraceful,
		signalChan:   make(chan os.Signal),
		shutdownChan: make(chan bool),
	}
}

func (srv *Server) ListenAndServe() error {
	addr := srv.HttpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.Listener = ln
	return srv.Serve()
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := srv.HttpServer.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.HttpServer.TLSConfig != nil {
		*config = *srv.HttpServer.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.Listener = tls.NewListener(ln, config)
	return srv.Serve()
}

func (srv *Server) Serve() error {
	go srv.handleSignals()
	err := srv.HttpServer.Serve(srv.Listener)

	srv.logf("waiting for connections closed.")
	<-srv.shutdownChan
	srv.logf("all connections closed.")

	return err
}

func (srv *Server) getNetListener(addr string) (net.Listener, error) {
	var ln net.Listener
	var err error

	if srv.isGraceful {
		file := os.NewFile(GRACEFUL_LISTENER_FD, "")
		ln, err = net.FileListener(file)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v %s", err, srv.HttpServer.Addr)
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return nil, err
		}
	}
	return ln, nil
}

func (srv *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)

	for {
		sig = <-srv.signalChan
		switch sig {
		case syscall.SIGTERM:
			srv.logf("received SIGTERM, graceful shutting down HTTP server.")
			srv.shutdownHTTPServer()
		case syscall.SIGUSR2:
			if srv.restartFlag {
				srv.logf("received SIGUSR2, graceful restarting %s HTTP server.", srv.HttpServer.Addr)
				if pid, err := srv.startNewProcess(); err != nil {
					srv.logf("start new process failed: %v, continue serving.", err)
				} else {
					srv.logf("start new process successed, the new pid is %d.", pid)
					srv.shutdownHTTPServer()
				}
			}
		default:
		}
	}
}

func (srv *Server) shutdownHTTPServer() {
	if err := srv.HttpServer.Shutdown(context.Background()); err != nil {
		srv.logf("HTTP server shutdown error: %v", err)
	} else {
		srv.logf("HTTP server shutdown success.")
		srv.shutdownChan <- true
	}
}

// start new process to handle HTTP Connection
func (srv *Server) startNewProcess() (uintptr, error) {
	listenerFd, err := srv.getTCPListenerFd()
	if err != nil { //
		return 0, fmt.Errorf("failed to get socket file descriptor: %v", err)
	}
	// set graceful restart env flag
	envs := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			envs = append(envs, value)
		}
	}

	envs = append(envs, GRACEFUL_ENVIRON_STRING)
	//fmt.Println("env:", envs)
	execSpec := &syscall.ProcAttr{
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}
	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return uintptr(fork), nil
}

func (srv *Server) getTCPListenerFd() (uintptr, error) {
	file, err := srv.Listener.(*net.TCPListener).File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}

func (srv *Server) getTCPListenerFdList() ([]uintptr, error) {
	var files = make([]uintptr, len(srv.LinsterList))
	for index, val := range srv.LinsterList {
		file, err := val.(*net.TCPListener).File()
		fmt.Println("val:", val, "file.fd:", file.Fd())
		if err != nil {
			return nil, err
		}
		files[index] = file.Fd()
	}
	return files, nil
}

func (srv *Server) logf(format string, args ...interface{}) {
	pids := strconv.Itoa(os.Getpid())
	format = "[pid " + pids + "] " + format

	if srv.HttpServer.ErrorLog != nil {
		srv.HttpServer.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
