package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/yungbote/slotter/backend/services/database/internal/events"
)

// Handler holds references to your subscriber and the Upgrader.
type Handler struct {
	Subscriber *events.PubSubSubscriber
	Upgrader   websocket.Upgrader
}

// NewHandler constructor. Typically called from main or your router setup.
func NewHandler(subscriber *events.PubSubSubscriber) *Handler {
	return &Handler{
		Subscriber: subscriber,
		Upgrader: websocket.Upgrader{
			// In production, do a stricter CheckOrigin.
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWSConnection is your GET /ws endpoint (or /ws/:company_id, etc.) that upgrades to WebSocket.
func (h *Handler) HandleWSConnection(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user_id in context"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	// If "company_id" is also in context, we capture it:
	var companyID uuid.UUID
	if cidVal, ok := c.Get("company_id"); ok {
		if cid, isUUID := cidVal.(uuid.UUID); isUUID {
			companyID = cid
		}
	}

	// Upgrade to WebSocket
	ws, err := h.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	// Subscribe to user channel
	userChan, err := h.Subscriber.SubscribeUser(userID)
	if err != nil {
		log.Printf("SubscribeUser error: %v", err)
		return
	}

	// Subscribe to company channel if we have a valid companyID
	var companyChan <-chan *redis.Message
	if companyID != uuid.Nil {
		companyChan, err = h.Subscriber.SubscribeCompany(companyID)
		if err != nil {
			log.Printf("SubscribeCompany error: %v", err)
			// Not fatal, we can continue with just userChan
		}
	}

	// Read from subscriptions in a separate goroutine
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)

		for {
			select {
			case msg, ok := <-userChan:
				if !ok {
					log.Println("User subscription channel closed.")
					return
				}
				// Forward the raw message payload to WebSocket
				errWrite := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if errWrite != nil {
					log.Printf("Error writing user message to WS: %v", errWrite)
					return
				}

			case msg, ok := <-companyChan:
				if companyChan == nil {
					continue
				}
				if !ok {
					log.Println("Company subscription channel closed.")
					return
				}
				// Forward the raw message
				errWrite := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if errWrite != nil {
					log.Printf("Error writing company message to WS: %v", errWrite)
					return
				}
			}
		}
	}()

	// Meanwhile, read from WS to keep the connection open (or handle client messages)
	for {
		_, _, errRead := ws.ReadMessage()
		if errRead != nil {
			log.Printf("WebSocket read error (userID=%s): %v", userID, errRead)
			break
		}
	}

	// On disconnect, unsubscribe
	h.Subscriber.UnsubscribeUser(userID)
	if companyID != uuid.Nil {
		h.Subscriber.UnsubscribeCompany(companyID)
	}

	<-doneChan
}

