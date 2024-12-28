package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/mcpnotion/tools"
)

const (
	serverName    = "mcpnotion"
	serverVersion = "0.0.1"
)

func main() {
	// get token from command line
	argsToken := flag.String("token", "", "notion token")
	argsLevel := flag.String("level", "info", "log level")
	argsDebug := flag.String("debug", "", "debug file")
	flag.Parse()

	// print all the args on stderr
	token := *argsToken
	level := *argsLevel
	debugFile := *argsDebug

	if token == "" {
		// retrive the token from the environment variable
		env_token := os.Getenv("NOTION_TOKEN")
		if env_token == "" {
			fmt.Println("Notion token is required")
			flag.PrintDefaults()
			os.Exit(1)
		}
		token = env_token
	}

	// create the mcpServerDefinition
	mcpServerDefinition := gomcp.NewMcpServerDefinition(serverName, serverVersion)
	mcpServerDefinition.SetDebugLevel(level, debugFile)

	// add the tools configuration to the server configuration
	mcpToolsDefinition := mcpServerDefinition.WithTools(&tools.NotionToolConfiguration{
		NotionToken: token,
	}, tools.NotionToolInit)

	mcpToolsDefinition.AddTool("notion_get_page", "Get the markdown content of a notion page", tools.NotionGetPage)
	mcpToolsDefinition.AddTool("ping", "A ping function", tools.NotionPing)

	mcp, err := gomcp.NewModelContextProtocolServer(mcpServerDefinition)
	if err != nil {
		fmt.Println("Error creating MCP server:", err)
		os.Exit(1)
	}
	// start the server
	transport := mcp.StdioTransport()
	err = mcp.Start(transport)
	if err != nil {
		fmt.Println("Error starting MCP server:", err)
		os.Exit(1)
	}
}
