package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}

// Run executable file from command line
// [out][port]
// start statis file server
// start websocket server
func run() error {
	fmt.Printf("Enter your port: ")
	var port string
	_, err := fmt.Scanln(&port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(l.Addr())

	cs := newChatServer()
	s := &http.Server{
		Handler:        cs,               //cs ServeHTTP
		ReadTimeout:    time.Second * 10, //10s
		WriteTimeout:   time.Second * 10, //10s
		MaxHeaderBytes: 1 << 20,          // 1 Mb
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case s := <-c:
		fmt.Printf("got signal to terminating executable %v\n", s)
	case e := <-errc:
		fmt.Printf("unexpected event happend %v \n", e)
	}

	// To handle grace full shut down server,
	// using context time out make sure context is in cancel state
	// refuse all other request . then after time out server shut down
	// discussion
	// https://www.reddit.com/r/golang/comments/16l1mhw/looking_for_clarifications_on_how_graceful/
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return s.Shutdown(ctx)
}
