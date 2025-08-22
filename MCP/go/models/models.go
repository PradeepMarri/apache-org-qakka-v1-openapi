package models

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Definition mcp.Tool
	Handler    func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// ApiResponse represents the ApiResponse schema from the OpenAPI specification
type ApiResponse struct {
	Count int `json:"count,omitempty"` // Count of Queues or QueueMessages returned by the call.
	Message string `json:"message,omitempty"` // Informative message intended for client.
	Queuemessages []QueueMessage `json:"queueMessages,omitempty"` // Queues Messages returned by the call, or empty if not applicable.
	Queues []Queue `json:"queues,omitempty"` // Queues returned but the call, or empty if not applicable.
}

// Queue represents the Queue schema from the OpenAPI specification
type Queue struct {
	Name string `json:"name"` // Name of queue, must be unique.
}

// QueueMessage represents the QueueMessage schema from the OpenAPI specification
type QueueMessage struct {
	Queuemessageid string `json:"queueMessageId,omitempty"` // UUID of Queue Message in local region.
	Sendingregion string `json:"sendingRegion,omitempty"` // Region from which was sent
	Href string `json:"href,omitempty"` // URL of data associated with Queue Message (if not embedded JSON)
	Receivingregion string `json:"receivingRegion,omitempty"` // Regions to which message will be sent
	Createdate int64 `json:"createDate,omitempty"` // Date that message was received by system.
	Queuename string `json:"queueName"` // Name of Queue for message.
	Data string `json:"data,omitempty"` // Embedded JSON to be sent with Queue Message.
	Messageid string `json:"messageId,omitempty"` // UUID of Message Data associated with this Queue Message
	Contenttype string `json:"contentType,omitempty"` // Content-type of data associated with QueueMessage.
}
