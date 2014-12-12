package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"syscall"
	"time"
)

func dialRemoteLogger(addr string) (io.Writer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Duration(5)*time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Status(message string) {
	fmt.Fprintf(os.Stderr, "reporter:status:%s\n", message)
	os.Stderr.Sync()
}

type prefixLogger struct {
	prefix []byte
	w      io.Writer
}

func (p *prefixLogger) Write(b []byte) (n int, err error) {
	n, err = p.w.Write(p.prefix)
	if err != nil {
		return n, err
	}
	nn, err := p.w.Write(b)
	return n + nn, err
}

// NewPrefixLogger returns a writer that behaves like w except
// that it writes a prefix before each write
func newPrefixLogger(prefix string, w io.Writer) io.Writer {
	return &prefixLogger{[]byte(prefix), w}
}

// InitSignal register signals handler.
func InitSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	return c
}

// HandleSignal fetch signal from chan then do exit or reload.
func HandleSignal(c chan os.Signal) {
	// Block until a signal is received.
	for {
		s := <-c
		fmt.Println("comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
			return
		}
	}
}

func main() {
	conn, err := dialRemoteLogger("127.0.0.1:64213")
	if err != nil {
		Status(fmt.Sprintf("error dialing remote logger %s", err))
	} else {
		hostname, _ := os.Hostname()
		w := newPrefixLogger(fmt.Sprintf("[%s %s:%d] ", hostname, "1", 2), conn)
		log.SetOutput(w)

	}
	log.Println("fdffdsfsd")
	log.Println("4324342432")

}
