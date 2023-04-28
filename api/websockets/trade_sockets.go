package BitgetWebSockets

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func Connect() error {

	//variables initialization
	apiKey := "your_api_key"
	passphrase := "your_passphrase"
	secretKey := "your_secret_key"

	url := "wss://ws.bitget.com/spot/v1/stream"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("Error connecting to WebSocket: %v", err)
	}
	defer c.Close()

	// Send a message to the WebSocket
	message := []byte(`{
    "op": "login",
    "args": [{
        "apiKey": "<api_key>",
        "passphrase": "<passphrase>",
        "timestamp": "<timestamp>",
        "sign": "<sign>"
     }]
    }`)
	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("Error sending message to WebSocket: %v", err)
	}

	// Read messages from the WebSocket
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("Error reading message from WebSocket: %v", err)
		}
		fmt.Println("Received message from WebSocket:", string(message))
	}
}
