// package ws
package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	storageRoom "match/internal/storage/room"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var mu sync.Mutex
var Connections = make(map[string]*websocket.Conn)

type WSHandler struct {
	roomStore *storageRoom.RoomStorage
}

func NewWSHandler(roomStore *storageRoom.RoomStorage) *WSHandler {
	return &WSHandler{
		roomStore: roomStore,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WSHandler) HandleWSUpgrade(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда на WebSocket:", err)
		return
	}

	Connections[userID] = conn
	log.Printf("User %s connected via WS\n", userID)

	go ws.readMessages(userID, conn)
}

func (ws *WSHandler) readMessages(userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		mu.Lock()
		delete(Connections, userID)
		mu.Unlock()
		log.Printf("User %s disconnected\n", userID)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("readMessages err:", err)
			return
		}
		var msg map[string]interface{}
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		msgType, _ := msg["type"].(string)

		switch msgType {
		case "ACCEPT_ROOM":
			roomID, _ := msg["room_id"].(string)
			if roomID == "" {
				log.Println("нет room_id в ACCEPT_ROOM")
				continue
			}
			err := ws.handleAcceptRoom(userID, roomID)
			if err != nil {
				log.Println("handleAcceptRoom error:", err)
			}
		default:
			log.Println("Неизвестный msgType:", msgType)
		}
	}
}

func (ws *WSHandler) handleAcceptRoom(userID string, roomID string) error {
	ctx := context.TODO()
	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	r, err := ws.roomStore.GetRoom(ctx, id)
	if err != nil {
		log.Println("handleAcceptRoom: не нашли комнату:", err)
		return err
	}

	if len(r.UserIDs) < 2 {
		return fmt.Errorf("room has not enough participants")
	}
	creatorID := r.UserIDs[0]
	friendID := r.UserIDs[1]

	if userID != friendID.Hex() {
		return fmt.Errorf("not a friend, cannot accept")
	}

	creatorStr := creatorID.Hex()
	mu.Lock()
	creatorConn, ok := Connections[creatorStr]
	mu.Unlock()
	if !ok {
		log.Println("Создатель не в сети (нет WS-соединения)")
		return fmt.Errorf("creator not connected")
	}

	msg := map[string]interface{}{
		"type":    "ROOM_ACCEPTED",
		"room_id": roomID,
	}
	err = creatorConn.WriteJSON(msg)
	if err != nil {
		log.Println("Ошибка при отправке ROOM_ACCEPTED создателю:", err)
		return err
	}

	log.Printf("Отправили ROOM_ACCEPTED создателю %s по комнате %s\n", creatorStr, roomID)
	return nil
}

func (ws *WSHandler) SendInviteMessage(userID, fromUserID, roomID string) error {
	conn, ok := Connections[userID]
	if !ok {
		return fmt.Errorf("user %s not connected via WS", userID)
	}
	invite := map[string]string{
		"type":      "INVITATION",
		"room_id":   roomID,
		"from_user": fromUserID,
		"message":   "Вас пригласили в комнату",
	}
	return conn.WriteJSON(invite)
}
