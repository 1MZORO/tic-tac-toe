package ws

type RoomMessage struct {
    Action   string `json:"action"`
    RoomID   string `json:"roomId"`
    Symbol   string `json:"symbol,omitempty"`  
    Position int    `json:"position,omitempty"`
}
