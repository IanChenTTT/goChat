package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Messages are send on the msgs channel
// cannot keep up the message, closeSlow is called
type sub struct {
    msgs  chan []byte
    closeSlow func()
}
type chatServer struct{
    subMsgBuf int // Max number of message that can queue

    pubLimiter *rate.Limiter

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
    err := cs.serverMux.HandleFunc("/sub", cs.subsHandler)

    cs.serverMux.HandleFunc("pub", cs.pubHandler)
    return cs
}
// Server each end point
func (cs *chatServer)ServeHTTP(w http.ResponseWriter, r *http.Request ){
    cs.serverMux.ServeHTTP(w,r)
}
func subsHandler(w http.ResponseWriter, r *http.Request ) error



