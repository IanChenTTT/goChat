package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// Messages are send on the msgs channel
// cannot keep up the message, closeSlow is called
type sub struct {
    msgs  chan []byte
    closeSlow func()
}
type chatServer struct{
    subMsgBuf int // Max number of message that can queue

    pubLimiter *rate.Limiter // connect limmiter

    logf func(f string, v ...interface{})

    serverMux http.ServeMux //For various end point

    subsMut sync.Mutex // Mutex for subs map

    subs map[*sub]struct{}
}
// return default chat server
// subMsgBuf : 16
// logf log.Printf TODO
// pubLimiter : 100 mil 8 token
// route / : ./
// route /sub
// route /pub
func newChatServer() *chatServer{
    cs := &chatServer{
        subMsgBuf: 16,
        logf: log.Printf,
        subs: make(map[*sub]struct{}),
        pubLimiter: rate.NewLimiter(rate.Every(time.Microsecond*100),8),
    }
    cs.serverMux.Handle("/", http.FileServer(http.Dir("."))) //TODO
    cs.serverMux.HandleFunc("/sub", cs.subsHandler)

    //cs.serverMux.HandleFunc("pub", cs.pubHandler)
    return cs
}
// Serve each end point
func (cs *chatServer)ServeHTTP(w http.ResponseWriter, r *http.Request ){
    cs.serverMux.ServeHTTP(w,r)
}
// subscribeHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (cs *chatServer) subsHandler(w http.ResponseWriter, r *http.Request ) {
    err := cs.subscribe(r.Context(), w, r)
    if errors.Is(err, context.Canceled){
        return 
    }
    if websocket.CloseStatus(err) == websocket.StatusNormalClosure|| websocket.CloseStatus(err) == websocket.StatusGoingAway{
       return 
    }
    if err != nil {
        cs.logf("%v", err)
        return 
    }
}

// addSubscriber register a subscriber
func (cs *chatServer) addSubscriber(s *sub){
    cs.subsMut.Lock()
    cs.subs[s] = struct{}{}
    cs.subsMut.Unlock()
}
// addSubscriber register a subscriber
func (cs *chatServer) deleteSubscriber(s *sub){
    cs.subsMut.Lock()
    delete(cs.subs, s)
    cs.subsMut.Unlock()
}
// if exceed timeout duration return err
func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error{
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    return c.Write(ctx, websocket.MessageText, msg)
}


// subscribe subscribes the given WebSocket to all broadcast messages.
// It creates a subscriber with a buffered msgs chan to give some room to slower
// connections and then registers the subscriber. It then listens for all messages
// and writes them to the WebSocket. If the context is cancelled or
// an error occurs, it returns and deletes the subscription.
//
// It uses CloseRead to keep reading from the connection to process control
// messages and cancel the context if the connection drops.
func (cs *chatServer) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	s := &sub{
		msgs: make(chan []byte, cs.subMsgBuf),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	cs.addSubscriber(s)
	defer cs.deleteSubscriber(s)

	c2, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
