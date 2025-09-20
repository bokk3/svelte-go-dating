package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"
    
    "github.com/gofiber/websocket/v2"
    "github.com/google/uuid"
    
    "dating-svelte/internal/database"
    "dating-svelte/internal/models"
)

type Hub struct {
    clients    map[uuid.UUID]*Client
    register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
    db         *database.DB
    mu         sync.RWMutex
}

type Client struct {
    id     uuid.UUID
    userID uuid.UUID
    conn   *websocket.Conn
    send   chan []byte
    hub    *Hub
}

type Message struct {
    Type      string      `json:"type"`
    MatchID   *uuid.UUID  `json:"match_id,omitempty"`
    Message   *string     `json:"message,omitempty"`
    UserID    *uuid.UUID  `json:"user_id,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data,omitempty"`
}

func NewHub(db *database.DB) *Hub {
    return &Hub{
        clients:    make(map[uuid.UUID]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte),
        db:         db,
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.userID] = client
            h.mu.Unlock()
            
            log.Printf("User %s connected", client.userID)
            
            // Send online status to matches
            h.broadcastUserStatus(client.userID, "online")
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.userID]; ok {
                delete(h.clients, client.userID)
                close(client.send)
            }
            h.mu.Unlock()
            
            log.Printf("User %s disconnected", client.userID)
            
            // Send offline status to matches
            h.broadcastUserStatus(client.userID, "offline")
            
        case message := <-h.broadcast:
            h.mu.RLock()
            for userID, client := range h.clients {
                select {
                case client.send <- message:
                default:
                    delete(h.clients, userID)
                    close(client.send)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) broadcastUserStatus(userID uuid.UUID, status string) {
    // Get user's matches
    matches, err := h.db.GetUserMatches(userID)
    if err != nil {
        return
    }
    
    statusMsg := Message{
        Type:      "user_status",
        UserID:    &userID,
        Timestamp: time.Now(),
        Data:      map[string]string{"status": status},
    }
    
    statusBytes, _ := json.Marshal(statusMsg)
    
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    for _, match := range matches {
        var targetUserID uuid.UUID
        if match.User1ID == userID {
            targetUserID = match.User2ID
        } else {
            targetUserID = match.User1ID
        }
        
        if client, ok := h.clients[targetUserID]; ok {
            select {
            case client.send <- statusBytes:
            default:
                delete(h.clients, targetUserID)
                close(client.send)
            }
        }
    }
}

func (h *Hub) SendMessageToMatch(matchID uuid.UUID, senderID uuid.UUID, message string) error {
    // Save message to database
    dbMessage := &models.Message{
        ID:          uuid.New(),
        MatchID:     matchID,
        SenderID:    senderID,
        Message:     message,
        MessageType: "text",
        CreatedAt:   time.Now(),
    }
    
    if err := h.db.CreateMessage(dbMessage); err != nil {
        return err
    }
    
    // Get match details to find recipient
    matches, err := h.db.GetUserMatches(senderID)
    if err != nil {
        return err
    }
    
    var recipientID uuid.UUID
    for _, match := range matches {
        if match.ID == matchID {
            if match.User1ID == senderID {
                recipientID = match.User2ID
            } else {
                recipientID = match.User1ID
            }
            break
        }
    }
    
    // Send to recipient if online
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    if client, ok := h.clients[recipientID]; ok {
        wsMessage := Message{
            Type:      "new_message",
            MatchID:   &matchID,
            Message:   &message,
            UserID:    &senderID,
            Timestamp: dbMessage.CreatedAt,
            Data:      dbMessage,
        }
        
        messageBytes, _ := json.Marshal(wsMessage)
        
        select {
        case client.send <- messageBytes:
        default:
            delete(h.clients, recipientID)
            close(client.send)
        }
    }
    
    return nil
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    for {
        _, messageBytes, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        
        var msg Message
        if err := json.Unmarshal(messageBytes, &msg); err != nil {
            continue
        }
        
        switch msg.Type {
        case "send_message":
            if msg.MatchID != nil && msg.Message != nil {
                c.hub.SendMessageToMatch(*msg.MatchID, c.userID, *msg.Message)
            }
        case "typing":
            if msg.MatchID != nil {
                // Handle typing indicators
                typingMsg := Message{
                    Type:      "typing",
                    MatchID:   msg.MatchID,
                    UserID:    &c.userID,
                    Timestamp: time.Now(),
                }
                typingBytes, _ := json.Marshal(typingMsg)
                c.hub.broadcast <- typingBytes
            }
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func HandleWebSocket(hub *Hub, userID uuid.UUID) func(*websocket.Conn) {
    return func(c *websocket.Conn) {
        client := &Client{
            id:     uuid.New(),
            userID: userID,
            conn:   c,
            send:   make(chan []byte, 256),
            hub:    hub,
        }
        
        client.hub.register <- client
        
        go client.writePump()
        go client.readPump()
    }
}