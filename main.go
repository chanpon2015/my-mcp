package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer("test", "v0.0.1", nil)
	//
	server.AddTools(&mcp.ServerTool{
		Tool: &mcp.Tool{
			Name:        "test",
			Description: "テスト用のツール",
		},
		Handler: testToolHandler,
	})
	server.AddPrompts(&mcp.ServerPrompt{
		Prompt: &mcp.Prompt{
			Arguments:   []*mcp.PromptArgument{},
			Description: "MCPのテスト",
			Name:        "test",
			Title:       "テスト",
		},
		Handler: testPromptHandler,
	})
	server.AddResources(&mcp.ServerResource{
		Resource: &mcp.Resource{
			Name:     "info",
			MIMEType: "text/plain",
			URI:      "embedded:info",
		},
		Handler: testResourceHandler,
	})

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)
	http.ListenAndServe(":8080", handler)
}

func testToolHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "テストです。"},
		},
	}, nil
}

func testPromptHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: "これはテストですと返答してください。",
				},
			},
		},
	}, nil
}

func testResourceHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.ReadResourceParams) (*mcp.ReadResourceResult, error) {
	u, err := url.Parse(params.URI)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "embedded" {
		return nil, fmt.Errorf("wrong scheme: %q", u.Scheme)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      params.URI,
				MIMEType: "text/plain",
				Text:     "テスト",
			},
		},
	}, nil
}
