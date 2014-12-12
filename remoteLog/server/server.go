package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// listen for log messages, and copy them to stderr
func startRemoteLogListner() string {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "0.0.0.0:0")
	if err != nil {
		log.Fatalf("error resolving %s", err)
	}
	ln, err := net.Listen("tcp4", tcpAddr.String())
	if err != nil {
		log.Fatalf("error listening %s", err)
	}
	log.Printf("listening on %v for log messages", ln.Addr())

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				continue
			}
			log.Printf("accepted remote logging connection from %s", conn.RemoteAddr())
			go io.Copy(os.Stderr, conn)
		}
	}()

	listenAddr := ln.Addr().String()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed getting hostname %s", err)
	}

	return strings.Replace(listenAddr, "0.0.0.0", hostname, 1)
}

func Status(message string) {
	fmt.Fprintf(os.Stderr, "reporter:status:%s\n", message)
	os.Stderr.Sync()
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
	loggerAddress := startRemoteLogListner()
	fmt.Println(loggerAddress)

	signalCH := InitSignal()
	HandleSignal(signalCH)

}
