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
	PageIdOrUrl string `json:"pageId" jsonschema_description:"the ID or the full URLof the Notion page to retrieve."`
}

// NotionGetDocumentContext is the context for the McpToolProcessor function.
type NotionGetDocumentContext struct {
	// The Notion client.
	NotionClient *notionapi.Client
}

// configuration for the McpToolProcessor function.
type NotionToolConfiguration struct {
	NotionToken string `json:"notionToken" jsonschema_description:"the notion token for the Notion client."`
}

// initializes the McpToolProcessor function.
func NotionToolInit(ctx context.Context, config *NotionToolConfiguration) (*NotionGetDocumentContext, error) {
	client := notionapi.NewClient(notionapi.Token(config.NotionToken))

	// we need to initialize the Notion client
	return &NotionGetDocumentContext{NotionClient: client}, nil
}

// retrieves the content of a Notion page identified by the PageId.
func NotionGetPage(ctx context.Context, toolCtx *NotionGetDocumentContext, input *NotionGetDocumentInput, output types.ToolCallResult) error {
	logger := gomcp.GetLogger(ctx)

	// extract the pageId from the input
	pageId := extractPageId(input.PageIdOrUrl)

	logger.Info("NotionGetPage", types.LogArg{
		"pageId": pageId,
	})

	content, err := getPageContent(ctx, toolCtx.NotionClient, pageId)
	if err != nil {
		return err
	}
	output.AddTextContent(strings.Join(content, "\n"))

	return nil
}

func extractPageId(pageIdOrUrl string) string {
	// if the pageIdOrUrl is a full URL, we extract the pageId from the URL
	// URL is like that:
	// https://www.notion.so/Test-Article-1696f674ce1e80bfbcdec283767f1395?pvs=4
	if strings.HasPrefix(pageIdOrUrl, "https://") {
		// we extract the pageId from the URL by isolating first the last segment of the path
		segments := strings.Split(pageIdOrUrl, "/")
		lastSegment := segments[len(segments)-1]
		// we remove everyting after the ?
		lastSegment = strings.Split(lastSegment, "?")[0]
		// we split by - and take the last one
		segments = strings.Split(lastSegment, "-")
		pageId := segments[len(segments)-1]
		return pageId
	}
	return pageIdOrUrl
}

type NotionPingInput struct {
	Message string `json:"message" jsonschema_description:"the message to ping."`
}

func NotionPing(ctx context.Context, toolCtx *NotionGetDocumentContext, input *NotionPingInput, output types.ToolCallResult) error {
	output.AddTextContent("pong: " + input.Message)
	return nil
}
