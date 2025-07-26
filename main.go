package main

import (
	"fmt"
	"net/http"

	"github.com/1MZORO/tiktactoe/ws"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWebSocket)
	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
