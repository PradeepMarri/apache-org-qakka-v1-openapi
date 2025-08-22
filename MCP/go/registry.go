package main

import (
	"github.com/qakka/mcp-server/config"
	"github.com/qakka/mcp-server/models"
	tools_queues "github.com/qakka/mcp-server/tools/queues"
	tools_status "github.com/qakka/mcp-server/tools/status"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_queues.CreateGetlistofqueuesTool(cfg),
		tools_queues.CreateCreatequeueTool(cfg),
		tools_queues.CreateDeletequeueTool(cfg),
		tools_queues.CreateGetqueueconfigTool(cfg),
		tools_queues.CreateUpdatequeueconfigTool(cfg),
		tools_queues.CreateGetmessagedataTool(cfg),
		tools_queues.CreateGetnextmessagesTool(cfg),
		tools_queues.CreateAckmessageTool(cfg),
		tools_status.CreateStatusTool(cfg),
	}
}
