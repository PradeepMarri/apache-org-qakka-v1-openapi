package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/qakka/mcp-server/config"
	"github.com/qakka/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetmessagedataHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queueNameVal, ok := args["queueName"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: queueName"), nil
		}
		queueName, ok := queueNameVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: queueName"), nil
		}
		queueMessageIdVal, ok := args["queueMessageId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: queueMessageId"), nil
		}
		queueMessageId, ok := queueMessageIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: queueMessageId"), nil
		}
		url := fmt.Sprintf("%s/queues/%s/data/%s", cfg.BaseURL, queueName, queueMessageId)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No authentication required for this endpoint
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.ApiResponse
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetmessagedataTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_queues_queueName_data_queueMessageId",
		mcp.WithDescription("Get data associated with a Queue Message."),
		mcp.WithString("queueName", mcp.Required(), mcp.Description("Name of Queue")),
		mcp.WithString("queueMessageId", mcp.Required(), mcp.Description("ID of Queue Message for which data is to be returned")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetmessagedataHandler(cfg),
	}
}
