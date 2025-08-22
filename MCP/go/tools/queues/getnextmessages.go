package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/qakka/mcp-server/config"
	"github.com/qakka/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetnextmessagesHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		queryParams := make([]string, 0)
		if val, ok := args["count"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("count=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/queues/%s/messages%s", cfg.BaseURL, queueName, queryString)
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

func CreateGetnextmessagesTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_queues_queueName_messages",
		mcp.WithDescription("Get next Queue Messages from a Queue"),
		mcp.WithString("queueName", mcp.Required(), mcp.Description("Name of Queue")),
		mcp.WithString("count", mcp.Description("Number of messages to get")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetnextmessagesHandler(cfg),
	}
}
