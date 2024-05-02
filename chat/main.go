// main is entry point of server side
// code base play ground, educational
// purpose, fell free to modify
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct{
  conns map[*websocket.Conn]bool //Websocket client map
  mu sync.RWMutex //Handle race condition in a connection map
}
func NewServer() *Server {
  return &Server{
    conns: make(map[*websocket.Conn]bool),
  }
}
// handleWS is a method from websocket server
// that is sync mutex for each incomming connection
func (s *Server) handleWS(ws *websocket.Conn){
  fmt.Println("new incomming from client", ws.RemoteAddr())
  s.mu.Lock()
  s.conns[ws] = true
  s.mu.Unlock()
  
  s.readLoop(ws)
}
func (s *Server) readLoop(ws *websocket.Conn){
  b := make([]byte, 1024)
  for{
    n, err := ws.Read(b)
    if err != nil{
      if err == io.EOF{
        break
      }
      fmt.Println("read error", err)
      continue
    }
    msg := b[:n]
    fmt.Println(string(msg))

    ws.Write([]byte("thank you for the msg"))
  }
}
func main() {
  serer := NewServer()
  http.Handle("/ws", websocket.Handler(serer.handleWS))
  http.ListenAndServe(":3000", nil)

}
