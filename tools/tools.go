package tools

import (
	"context"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/gomcp/types"
)

// describes the page we want to retrieve.
type NotionGetDocumentInput struct {
	PageId string `json:"pageId" jsonschema_description:"the ID of the Notion page to retrieve."`
}

// NotionGetDocumentContext is the context for the McpToolProcessor function.
type NotionGetDocumentContext struct {
	// The Notion client.
	NotionClient *notionapi.Client
}

// configuration for the McpToolProcessor function.
type NotionGetDocumentConfiguration struct {
	NotionToken string `json:"notionToken" jsonschema_description:"the notion token for the Notion client."`
}

// initializes the McpToolProcessor function.
func NotionToolInit(ctx context.Context, config *NotionGetDocumentConfiguration) (*NotionGetDocumentContext, error) {
	client := notionapi.NewClient(notionapi.Token(config.NotionToken))

	// we need to initialize the Notion client
	return &NotionGetDocumentContext{NotionClient: client}, nil
}

// retrieves the content of a Notion page identified by the PageId.
func NotionGetPage(ctx context.Context, toolCtx *NotionGetDocumentContext, input *NotionGetDocumentInput, output types.ToolCallResult) error {
	logger := gomcp.GetLogger(ctx)
	logger.Info("NotionGetPage", types.LogArg{
		"pageId": input.PageId,
	})

	content, err := getPageContent(ctx, toolCtx.NotionClient, input.PageId)
	if err != nil {
		return err
	}
	output.AddTextContent(strings.Join(content, "\n"))

	return nil
}

func RegisterTools(toolRegistry types.ToolRegistry) error {
	toolProvider, err := toolRegistry.DeclareToolProvider("notion", NotionToolInit)
	if err != nil {
		return err
	}
	err = toolProvider.AddTool("notion_get_page", "Get the markdown content of a notion page", NotionGetPage)
	if err != nil {
		return err
	}
	return nil
}
