package coinbase

const DefaultURL = "wss://ws-feed.exchange.coinbase.com"

// RequestType represents the type of request.
type RequestType string

const (
	RequestTypeSubscribe   RequestType = "subscribe"
	RequestTypeUnsubscribe RequestType = "unsubscribe"
	RequestTypeError       RequestType = "error"
)

// ChannelType represents the type of channel on Coinbase.
type ChannelType string

const ChannelTypeMatches ChannelType = "matches"

type Channel struct {
	Name       ChannelType
	ProductIDs []string
}

// Request is a request to be sent to the Coinbase websocket.
type Request struct {
	Type       RequestType `json:"type"`
	ProductIDs []string    `json:"product_ids"`
	Channels   []Channel   `json:"channels"`
}

// Response is the response received after a request submission.
type Response struct {
	Type      string    `json:"type"`
	Size      string    `json:"size"`
	Price     string    `json:"price"`
	ProductID string    `json:"product_id"`
	Channels  []Channel `json:"channels"`
	Message   string    `json:"message,omitempty"`
}
