package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/mcpnotion/tools"
)

func main() {
	configFile := flag.String("configFile", "", "config file path (required)")
	flag.Parse()

	if *configFile == "" {
		fmt.Println("Config file is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	mcp, err := gomcp.NewModelContextProtocolServer(*configFile)
	if err != nil {
		fmt.Println("Error creating MCP server:", err)
		os.Exit(1)
	}
	toolRegistry := mcp.GetToolRegistry()

	err = tools.RegisterTools(toolRegistry)
	if err != nil {
		fmt.Println("Error registering tools:", err)
		os.Exit(1)
	}

	transport := mcp.StdioTransport()

	mcp.Start("mcpnotion", "0.1.0", transport)
}
