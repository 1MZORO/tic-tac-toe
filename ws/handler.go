package ws

import (
	"fmt"
	"log"
	"net/http"	
	"github.com/1MZORO/tiktactoe/game"
	"github.com/1MZORO/tiktactoe/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	fmt.Println("Client connected")

	for {
		var msg RoomMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			conn.Close()
			return
		}

		log.Printf("📨 Received: %+v", msg)

		switch msg.Action {
		case "create":
			if _, exists := models.Rooms[msg.RoomID]; exists {
				conn.WriteJSON(map[string]string{"error": "Room already exists"})
				continue
			}
			player := &models.Player{Conn: conn, ID: "id-x01"}

			room := &models.Room{
				ID:      msg.RoomID,
				PlayerX: player,
				Game:    &game.Game{Turn: "X"},
			}

			models.Rooms[msg.RoomID] = room
			conn.WriteJSON(map[string]string{"status": "waiting", "message": "Room created. Waiting for player O"})

		case "join":
			room, exists := models.Rooms[msg.RoomID]
			if !exists {
				conn.WriteJSON(map[string]string{"error": "Room does not exist"})
				continue
			}
			if room.PlayerO != nil {
				conn.WriteJSON(map[string]string{"error": "Room is full"})
				continue
			}
			player := &models.Player{Conn: conn, ID: "id-x02"}
			room.PlayerO = player
			room.PlayerX.Conn.WriteJSON(map[string]string{"status": "ready", "yourSymbol": "X"})
			room.PlayerO.Conn.WriteJSON(map[string]string{"status": "ready", "yourSymbol": "O"})

		case "move":
			handleMove(msg, conn) // You'll define this logic to handle game moves

		default:
			conn.WriteJSON(map[string]string{"error": "Invalid action"})
		}
	}
}

func handleMove(msg RoomMessage, conn *websocket.Conn) {
	row := msg.Position / 3
	col := msg.Position % 3
	room, exists := models.Rooms[msg.RoomID]
	if !exists {
		conn.WriteJSON(map[string]string{"error": "Room not found"})
		return
	}
	if msg.Symbol != "X" && msg.Symbol != "O" {
		conn.WriteJSON(map[string]string{"error": "Invalid symbol"})
		return
	}

	// Ensure player is in the room
	if (msg.Symbol == "X" && room.PlayerX.Conn != conn) ||
		(msg.Symbol == "O" && room.PlayerO.Conn != conn) {
		conn.WriteJSON(map[string]string{"error": "You are not this player"})
		return
	}

	// Ensure it's their turn
	if msg.Symbol != room.Game.Turn {
		conn.WriteJSON(map[string]string{"error": "Not your turn"})
		return
	}

	// Validate position
	if msg.Position < 0 || msg.Position >= 9 {
		conn.WriteJSON(map[string]string{"error": "Invalid board position"})
		return
	}

	// Check if already filled
	if room.Game.Board[row][col] != "" {
		conn.WriteJSON(map[string]string{"error": "Position already taken"})
		return
	}
	room.Game.Board[row][col] = msg.Symbol

	// Increment move count
	room.Game.Count++

	// Print Board
	room.Game.PrintBoard()

	// Check for winner
	if room.Game.CheckWinner() && room.Game.Count >=5{
		fmt.Printf("Called ! %d \n",room.Game.Count)
		winUpdate := map[string]interface{}{
			"action": "gameOver",
			"winner": room.Game.Winner,
			"board": room.Game.Board,
		}
		room.PlayerX.Conn.WriteJSON(winUpdate)
		room.PlayerO.Conn.WriteJSON(winUpdate)
		room.Game.Finished = true
		fmt.Printf("Game Over! Winner: %s\n", room.Game.Winner)
		return
	}

	//  Check for draw
	if room.Game.CheckDraw() {
		drawUpdate := map[string]interface{}{
			"action": "draw",
			"board":  room.Game.Board,
		}
		room.PlayerX.Conn.WriteJSON(drawUpdate)
		room.PlayerO.Conn.WriteJSON(drawUpdate)
		return
	}

	// Switch turn

	if room.Game.Turn == "X" {
		room.Game.Turn = "O"
	} else {
		room.Game.Turn = "X"
	}

	// Broadcast the updated board and next turn
	update := map[string]interface{}{
		"action":   "update",
		"board":    room.Game.Board,
		"nextTurn": room.Game.Turn,
	}
	room.PlayerX.Conn.WriteJSON(update)
	room.PlayerO.Conn.WriteJSON(update)
} 
