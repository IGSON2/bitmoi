package api

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func (s *Server) websocketTest(c *websocket.Conn) {
	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		err = c.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
