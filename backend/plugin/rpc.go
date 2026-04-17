package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

type JsonRpcMessage struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *int64          `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JsonRpcError   `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RequestHandler func(params json.RawMessage) (any, error)

type JsonRpcConn struct {
	mu        sync.Mutex
	writer    io.Writer
	scanner   *bufio.Scanner
	nextID    atomic.Int64
	pending   map[int64]chan *JsonRpcMessage
	handlers  map[string]RequestHandler
	done      chan struct{}
	closeOnce sync.Once
}

// NewJsonRpcConn creates a new JSON-RPC connection over the given reader/writer pair.
// Call Start() to begin reading messages from the reader.
func NewJsonRpcConn(reader io.Reader, writer io.Writer) *JsonRpcConn {
	scanner := bufio.NewScanner(reader)
	// Use a 10MB buffer to accommodate large tool results.
	scanner.Buffer(make([]byte, 0, 10*1024*1024), 10*1024*1024)

	return &JsonRpcConn{
		writer:   writer,
		scanner:  scanner,
		pending:  make(map[int64]chan *JsonRpcMessage),
		handlers: make(map[string]RequestHandler),
		done:     make(chan struct{}),
	}
}

// Start begins the background read loop that processes incoming JSON-RPC messages.
func (c *JsonRpcConn) Start() {
	go c.readLoop()
}

// Close signals the read loop to stop. It is safe to call multiple times.
func (c *JsonRpcConn) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
}

// Call sends a JSON-RPC request and blocks until a response is received or a 30-second timeout elapses.
func (c *JsonRpcConn) Call(method string, params any) (json.RawMessage, error) {
	rawParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}

	id := c.nextID.Add(1)
	ch := make(chan *JsonRpcMessage, 1)

	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	msg := &JsonRpcMessage{
		Jsonrpc: "2.0",
		ID:      &id,
		Method:  method,
		Params:  rawParams,
	}

	if err := c.send(msg); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("send request: %w", err)
	}

	select {
	case resp := <-ch:
		if resp.Error != nil {
			return nil, fmt.Errorf("rpc error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	case <-time.After(30 * time.Second):
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("rpc call %s timed out after 30s", method)
	case <-c.done:
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("connection closed")
	}
}

// Notify sends a JSON-RPC notification (a message with no ID that expects no response).
func (c *JsonRpcConn) Notify(method string, params any) error {
	rawParams, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal params: %w", err)
	}

	msg := &JsonRpcMessage{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  rawParams,
	}

	return c.send(msg)
}

// RegisterHandler registers a handler function for a given JSON-RPC method name.
func (c *JsonRpcConn) RegisterHandler(method string, handler RequestHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[method] = handler
}

// readLoop continuously reads newline-delimited JSON messages from the scanner.
func (c *JsonRpcConn) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("jsonrpc readLoop panic", fmt.Errorf("%v", r))
		}
	}()

	for {
		select {
		case <-c.done:
			return
		default:
		}

		if !c.scanner.Scan() {
			if err := c.scanner.Err(); err != nil {
				select {
				case <-c.done:
					return
				default:
					logger.Error("jsonrpc scanner error", err)
				}
			}
			return
		}

		line := c.scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var msg JsonRpcMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			logger.Warm("jsonrpc unmarshal error: %v, line: %s", err, string(line))
			continue
		}

		if msg.ID != nil && msg.Method == "" {
			// This is a response to a pending request.
			c.mu.Lock()
			ch, ok := c.pending[*msg.ID]
			if ok {
				delete(c.pending, *msg.ID)
			}
			c.mu.Unlock()

			if ok {
				ch <- &msg
			} else {
				logger.Warm("jsonrpc received response for unknown id: %d", *msg.ID)
			}
		} else if msg.Method != "" {
			// This is an incoming request or notification.
			c.mu.Lock()
			handler, ok := c.handlers[msg.Method]
			c.mu.Unlock()

			if !ok {
				logger.Warm("jsonrpc no handler for method: %s", msg.Method)
				if msg.ID != nil {
					resp := &JsonRpcMessage{
						Jsonrpc: "2.0",
						ID:      msg.ID,
						Error: &JsonRpcError{
							Code:    -32601,
							Message: "method not found: " + msg.Method,
						},
					}
					if err := c.send(resp); err != nil {
						logger.Error("jsonrpc send error response", err)
					}
				}
				continue
			}

			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Error("jsonrpc handler panic", fmt.Errorf("%v", r))
						if msg.ID != nil {
							resp := &JsonRpcMessage{
								Jsonrpc: "2.0",
								ID:      msg.ID,
								Error: &JsonRpcError{
									Code:    -32603,
									Message: "internal error",
								},
							}
							_ = c.send(resp)
						}
					}
				}()

				result, err := handler(msg.Params)
				if msg.ID != nil {
					resp := &JsonRpcMessage{
						Jsonrpc: "2.0",
						ID:      msg.ID,
					}
					if err != nil {
						resp.Error = &JsonRpcError{
							Code:    -32603,
							Message: err.Error(),
						}
					} else {
						rawResult, marshalErr := json.Marshal(result)
						if marshalErr != nil {
							resp.Error = &JsonRpcError{
								Code:    -32603,
								Message: "failed to marshal result: " + marshalErr.Error(),
							}
						} else {
							resp.Result = rawResult
						}
					}
					if sendErr := c.send(resp); sendErr != nil {
						logger.Error("jsonrpc send response", sendErr)
					}
				}
			}()
		}
	}
}

// send marshals the message to JSON and writes it as a single line to the writer.
func (c *JsonRpcConn) send(msg *JsonRpcMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	data = append(data, '\n')
	_, err = c.writer.Write(data)
	return err
}
