package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// Client represents an LSP client connected to gopls
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	requestID  int64
	requests   map[any]chan *JSONRPCResponse
	requestsMu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	initialized bool
	rootURI     uri.URI
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewClient creates a new LSP client connected to gopls
func NewClient(rootPath string) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "gopls", "serve")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start gopls: %w", err)
	}

	// Convert path to URI using the standard library
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	rootURI := uri.File(absPath)

	client := &Client{
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		requestID: 0,
		requests:  make(map[any]chan *JSONRPCResponse),
		ctx:       ctx,
		cancel:    cancel,
		rootURI:   rootURI,
	}

	// Start reading responses
	go client.readResponses()

	return client, nil
}

// Initialize performs LSP initialization
func (c *Client) Initialize() error {
	if c.initialized {
		return nil
	}

	params := &protocol.InitializeParams{
		ProcessID: int32(os.Getpid()),
		Capabilities: protocol.ClientCapabilities{
			Workspace: &protocol.WorkspaceClientCapabilities{
				WorkspaceEdit: &protocol.WorkspaceClientCapabilitiesWorkspaceEdit{
					DocumentChanges: true,
				},
			},
			TextDocument: &protocol.TextDocumentClientCapabilities{
				Rename: &protocol.RenameClientCapabilities{
					DynamicRegistration: false,
					PrepareSupport:      false,
				},
			},
		},
		WorkspaceFolders: []protocol.WorkspaceFolder{
			{
				URI:  string(c.rootURI),
				Name: filepath.Base(string(c.rootURI)),
			},
		},
	}

	var result protocol.InitializeResult
	if err := c.sendRequest("initialize", params, &result); err != nil {
		return fmt.Errorf("initialize request failed: %w", err)
	}

	// Send initialized notification
	if err := c.sendNotification("initialized", &protocol.InitializedParams{}); err != nil {
		return fmt.Errorf("initialized notification failed: %w", err)
	}

	// Add small delay to let gopls process the notification
	time.Sleep(100 * time.Millisecond)

	c.initialized = true
	return nil
}

// Rename performs a rename operation
func (c *Client) Rename(filePath string, line, character int, newName string) (*protocol.WorkspaceEdit, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	// Convert file path to URI using the standard library
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fileURI := uri.File(absPath)

	// Create a custom params structure that matches what gopls expects
	params := map[string]any{
		"textDocument": map[string]any{
			"uri": string(fileURI),
		},
		"position": map[string]any{
			"line":      line,
			"character": character,
		},
		"newName": newName,
	}

	var result protocol.WorkspaceEdit
	if err := c.sendRequest("textDocument/rename", params, &result); err != nil {
		return nil, fmt.Errorf("rename request failed: %w", err)
	}

	return &result, nil
}

// Close closes the LSP client
func (c *Client) Close() error {
	if c.cancel != nil {
		c.cancel()
	}

	if c.stdin != nil {
		_ = c.stdin.Close()
	}

	if c.cmd != nil {
		_ = c.cmd.Wait()
	}

	return nil
}

// sendRequest sends a JSON-RPC request and waits for response
func (c *Client) sendRequest(method string, params any, result any) error {
	id := atomic.AddInt64(&c.requestID, 1)

	// Create request as a map to ensure proper JSON serialization
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}

	// Create response channel
	responseChan := make(chan *JSONRPCResponse, 1)
	c.requestsMu.Lock()
	c.requests[id] = responseChan
	c.requestsMu.Unlock()

	// Clean up channel when done
	defer func() {
		c.requestsMu.Lock()
		delete(c.requests, id)
		c.requestsMu.Unlock()
	}()

	// Send request
	if err := c.sendMessage(request); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for response with timeout
	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()

	select {
	case response := <-responseChan:
		if response.Error != nil {
			return fmt.Errorf("LSP error: %s", response.Error.Message)
		}

		if result != nil && response.Result != nil {
			if err := json.Unmarshal(response.Result, result); err != nil {
				return fmt.Errorf("failed to unmarshal result: %w", err)
			}
		}

		return nil
	case <-timeout.C:
		return fmt.Errorf("request timeout after 30 seconds for method %s", method)
	case <-c.ctx.Done():
		return fmt.Errorf("context cancelled")
	}
}

// sendNotification sends a JSON-RPC notification
func (c *Client) sendNotification(method string, params any) error {
	// Create notification as a map
	notification := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	return c.sendMessage(notification)
}

// sendMessage sends a message over the LSP connection
func (c *Client) sendMessage(message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	content := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)

	if _, err := c.stdin.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// readResponses reads responses from the LSP server
func (c *Client) readResponses() {
	reader := bufio.NewReader(c.stdout)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		// Read Content-Length header
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		if !strings.HasPrefix(line, "Content-Length:") {
			continue
		}

		// Parse content length
		lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
		contentLength, err := strconv.Atoi(lengthStr)
		if err != nil {
			continue
		}

		// Read empty line
		if _, err := reader.ReadString('\n'); err != nil {
			continue
		}

		// Read message content
		content := make([]byte, contentLength)
		if _, err := io.ReadFull(reader, content); err != nil {
			continue
		}

		// Parse JSON-RPC response
		var response JSONRPCResponse
		if err := json.Unmarshal(content, &response); err != nil {
			continue
		}

		// Handle response
		c.handleResponse(&response)
	}
}

// handleResponse handles a JSON-RPC response
func (c *Client) handleResponse(response *JSONRPCResponse) {
	// Normalize ID type - JSON unmarshaling might convert int64 to float64
	normalizedID := response.ID
	if f, ok := response.ID.(float64); ok && f == float64(int64(f)) {
		normalizedID = int64(f)
	}

	c.requestsMu.RLock()
	responseChan, exists := c.requests[normalizedID]
	c.requestsMu.RUnlock()

	if !exists {
		return
	}

	select {
	case responseChan <- response:
	default:
		// Channel full, ignore
	}
}
