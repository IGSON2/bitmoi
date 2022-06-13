package apiWebsocket

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// )

// type MessageCh struct {
// 	ch   chan []byte
// 	conn *websocket.Conn
// 	m    sync.Mutex
// }

// func Upgrade(c *fiber.Ctx) error {
// 	if websocket.IsWebSocketUpgrade(c) {
// 		c.Locals("allowed", true)
// 		fmt.Println("Websocket upgrade is complete.")
// 		return c.Next()
// 	}
// 	return fiber.ErrUpgradeRequired
// }

// func (m *MessageCh) close() {
// 	m.conn.Close()
// }

// func InitSocket(c *websocket.Conn) {
// 	mCH := &MessageCh{ch: make(chan []byte), conn: c}
// 	mCH.Read()
// 	// mCH.Write()
// }

// func (m *MessageCh) Read() {
// 	defer m.close()
// 	for {
// 		mt, msg, err := m.conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Connection Error ! : ", err)
// 			break
// 		}
// 		log.Printf("recv : %s Type : %d", msg, mt)

// 		var unMashaledMsg Message
// 		unmarshalErr := json.Unmarshal(msg, &unMashaledMsg)
// 		if unmarshalErr != nil {
// 			log.Panicln("err : ", unmarshalErr)
// 		}
// 		// HandleMessage(unMashaledMsg, m)
// 	}
// }

// func (m *MessageCh) Write() {
// 	defer m.close()
// 	for {
// 		// var readMsg = <-m.ch
// 		readMsg2 := Message{"Temp"}
// 		if err := m.conn.WriteJSON(readMsg2); err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// }
